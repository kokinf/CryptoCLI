package modules

const BlockSizeIDEA = 8 // IDEA работает с блоками по 64 бита (8 байт)
const KeySizeIDEA = 16  // IDEA-128 использует ключ из 128 бит (16 байт)
const RoundsIDEA = 8    // IDEA использует 8 раундов + 4 в последнем

// Разбиение 64-битного блока на 4 слова по 16 бит
func splitBlock(block []byte) [4]uint16 {
	var words [4]uint16
	for i := 0; i < 4; i++ {
		words[i] = uint16(block[2*i])<<8 | uint16(block[2*i+1])
	}
	return words
}

// Сборка 4 слов по 16 бит в 64-битный блок
func combineWords(words [4]uint16) []byte {
	block := make([]byte, BlockSizeIDEA)
	for i := 0; i < 4; i++ {
		block[2*i] = byte(words[i] >> 8)
		block[2*i+1] = byte(words[i])
	}
	return block
}

// mul — умножение в поле GF(2^16 + 1), модуль 0x1100B
// mul — умножение в поле GF(2^16 + 1)
func mul(a, b uint16) uint16 {
	if a == 0 {
		return 0x0001 - b // вместо 0x10001 → избегаем переполнения
	}
	if b == 0 {
		return 0x0001 - a
	}

	result := uint32(a) * uint32(b)
	low := result & 0xFFFF
	high := result >> 16

	return uint16(low ^ high)
}

// invMul — обратное умножение в поле GF(2^16 + 1)
func invMul(a uint16) uint16 {
	for b := uint16(1); b != 0; b++ {
		if mul(a, b) == 1 {
			return b
		}
	}
	return 0
}

// add — сложение по модулю 2^16
func add(a, b uint16) uint16 {
	return a + b
}

// invAdd — обратное сложение по модулю 2^16
func invAdd(a uint16) uint16 {
	return 0x0000 - a
}

// generateRoundKeys — генерация подключей IDEA
func generateRoundKeys(key []byte) [][6]uint16 {
	keySchedule := make([][6]uint16, RoundsIDEA+1)
	for r := 0; r <= RoundsIDEA; r++ {
		for j := 0; j < 6; j++ {
			keySchedule[r][j] = uint16(key[(r*6+j)*2])<<8 | uint16(key[(r*6+j)*2+1])
		}
	}
	return keySchedule
}

// ideaRound — один раунд IDEA
func ideaRound(data [4]uint16, keys [6]uint16) [4]uint16 {
	X1, X2, X3, X4 := data[0], data[1], data[2], data[3]
	K1, K2, K3, K4, K5, K6 := keys[0], keys[1], keys[2], keys[3], keys[4], keys[5]

	t1 := mul(X1, K1)
	t2 := add(X2, K2)
	t3 := add(X3, K3)
	t4 := mul(X4, K4)

	t5 := t1 ^ t3
	t6 := t2 ^ t4

	t7 := mul(t5, K5)
	t8 := add(t6, t7)
	t9 := add(t5, t8)

	Y1 := add(X3, t8)
	Y2 := add(X1, t9)
	Y3 := mul(X2, t9)
	Y4 := mul(X4, K6)

	return [4]uint16{Y1, Y2, Y3, Y4}
}

// IDEA_EncryptBlock — шифрует один 64-битный блок IDEA
func IDEA_EncryptBlock(block, key []byte) []byte {
	if len(block) != BlockSizeIDEA {
		panic("Неверная длина блока IDEA")
	}
	if len(key) != KeySizeIDEA {
		panic("Неверная длина ключа IDEA")
	}

	state := splitBlock(block)
	keySchedule := generateRoundKeys(key)

	for round := 0; round < RoundsIDEA; round++ {
		var keys [6]uint16
		copy(keys[:], keySchedule[round][:])
		state = ideaRound(state, keys)
	}

	// Последний раунд без перестановки
	final := state
	final[0], final[1], final[2], final[3] =
		mul(final[2], keySchedule[RoundsIDEA][0]),
		add(final[0], keySchedule[RoundsIDEA][1]),
		add(final[1], keySchedule[RoundsIDEA][2]),
		mul(final[3], keySchedule[RoundsIDEA][3])

	return combineWords(final)
}

// IDEA_DecryptBlock — расшифровывает один 64-битный блок IDEA
func IDEA_DecryptBlock(block, key []byte) []byte {
	if len(block) != BlockSizeIDEA {
		panic("Неверная длина блока IDEA")
	}
	if len(key) != KeySizeIDEA {
		panic("Неверная длина ключа IDEA")
	}

	state := splitBlock(block)
	keySchedule := generateRoundKeys(key)

	// Используем обратные раунды
	last := RoundsIDEA - 1
	keys := [6]uint16{
		invMul(keySchedule[last][5]), // K49⁻¹
		keySchedule[last][4],         // K50
		keySchedule[last][3],         // K51
		invMul(keySchedule[last][2]), // K52⁻¹
		0, 0,
	}
	state = ideaRound(state, keys)

	for r := last - 1; r >= 0; r-- {
		keys[0] = invMul(keySchedule[r][5]) // K_{r*6+5}⁻¹
		keys[1] = keySchedule[r][4]         // K_{r*6+4}
		keys[2] = keySchedule[r][3]         // K_{r*6+3}
		keys[3] = invMul(keySchedule[r][2]) // K_{r*6+2}⁻¹
		keys[4] = keySchedule[r][1]         // K_{r*6+1}
		keys[5] = keySchedule[r][0]         // K_{r*6+0}
		state = ideaRound(state, keys)
	}

	// Первый раунд без InvMixColumns
	keys = [6]uint16{
		invMul(keySchedule[0][5]), // K5⁻¹
		keySchedule[0][4],         // K6
		keySchedule[0][3],         // K7
		invMul(keySchedule[0][2]), // K8⁻¹
		0, 0,
	}
	state = ideaRound(state, keys)

	return combineWords(state)
}

// substitutionIDEA — шифрует несколько блоков IDEA
func substitutionIDEA(data []byte, ideaKey []byte) []byte {
	if len(ideaKey) != KeySizeIDEA {
		panic("Неверная длина ключа IDEA")
	}

	if len(data)%BlockSizeIDEA != 0 {
		data = pkcs7Padding(data, BlockSizeIDEA)
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeIDEA {
		block := data[i : i+BlockSizeIDEA]
		encrypted := IDEA_EncryptBlock(block, ideaKey)
		copy(result[i:], encrypted)
	}
	return result
}

// reverseSubstitutionIDEA — расшифровывает несколько блоков IDEA
func reverseSubstitutionIDEA(data []byte, ideaKey []byte) []byte {
	if len(ideaKey) != KeySizeIDEA {
		panic("Неверная длина ключа IDEA")
	}

	if len(data)%BlockSizeIDEA != 0 {
		data = pkcs7Padding(data, BlockSizeIDEA)
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeIDEA {
		block := data[i : i+BlockSizeIDEA]
		decrypted := IDEA_DecryptBlock(block, ideaKey)
		copy(result[i:], decrypted)
	}
	return result
}
