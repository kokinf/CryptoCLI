package modules

const BlockSizeIDEA = 8
const KeySizeIDEA = 16
const RoundsIDEA = 8

func init() {
	for b := 0; b < 256; b++ {
		mul9[b] = mul2(mul2(mul2(byte(b))^byte(b))) ^ byte(b)
		mul11[b] = mul2(mul2(mul2(byte(b)^mul2(byte(b)))) ^ mul2(byte(b)))
		mul13[b] = mul2(mul2(mul2(byte(b)^byte(b))) ^ mul2(byte(b)) ^ byte(b))
		mul14[b] = mul2(mul2(mul2(byte(b)^byte(b))) ^ byte(b))
	}
}

// invMul — обратное умножение в поле Галуа
func invMul(a uint16) uint16 {
	for b := uint16(1); b != 0; b++ {
		if (a*b)&0xff == 1 {
			return b
		}
	}
	return 0
}

// invAdd — обратное сложение
func invAdd(a uint16) uint16 {
	return 0 - a
}

// Сборка 4 слов по 16 бит в 64-битный блок
func combineWordsIDEA(words [4]uint16) []byte {
	block := make([]byte, BlockSizeIDEA)
	for i := 0; i < 4; i++ {
		block[2*i] = byte(words[i] >> 8)
		block[2*i+1] = byte(words[i])
	}
	return block
}

// generateRoundKeys — генерация подключей IDEA
func generateRoundKeys(key []byte) [][6]uint16 {
	keySchedule := make([][6]uint16, RoundsIDEA+1)
	for r := 0; r <= RoundsIDEA; r++ {
		for j := 0; j < 6; j++ {
			keySchedule[r][j] = uint16(key[(r*6+j)%KeySizeIDEA])<<8 |
				uint16(key[(r*6+j+1)%KeySizeIDEA])
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

func IDEA_EncryptBlock(block, key []byte) []byte {
	if len(block) != BlockSizeIDEA {
		panic("❌ Неверная длина блока IDEA")
	}
	if len(key) != KeySizeIDEA {
		panic("❌ Неверная длина ключа IDEA")
	}

	state := splitBlock(block)
	keySchedule := generateRoundKeys(key)

	for round := 0; round < RoundsIDEA; round++ {
		keys := keySchedule[round]
		state = ideaRound(state, keys)
	}

	// Последний раунд без перестановки
	final := state
	final[0], final[1], final[2], final[3] =
		mul(final[2], keySchedule[RoundsIDEA][0]),
		add(final[0], keySchedule[RoundsIDEA][1]),
		add(final[1], keySchedule[RoundsIDEA][2]),
		mul(final[3], keySchedule[RoundsIDEA][3])

	return combineWordsIDEA(final)
}

func IDEA_DecryptBlock(block, key []byte) []byte {
	if len(block) != BlockSizeIDEA {
		panic("❌ Неверная длина блока IDEA")
	}
	if len(key) != KeySizeIDEA {
		panic("❌ Неверная длина ключа IDEA")
	}

	state := splitBlock(block)
	keySchedule := generateRoundKeys(key)

	// Используем обратные ключи
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
		keys[0] = invMul(keySchedule[r][5])
		keys[1] = keySchedule[r][4]
		keys[2] = keySchedule[r][3]
		keys[3] = invMul(keySchedule[r][2])
		keys[4] = keySchedule[r][1]
		keys[5] = keySchedule[r][0]
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

	return combineWordsIDEA(state)
}

func substitutionIDEA(data []byte, ideaKey []byte) []byte {
	if len(ideaKey) != KeySizeIDEA {
		panic("❌ Ключ должен быть 16 байт")
	}

	data = pkcs7Padding(data, BlockSizeIDEA)

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeIDEA {
		block := data[i : i+BlockSizeIDEA]
		encrypted := IDEA_EncryptBlock(block, ideaKey)
		copy(result[i:], encrypted)
	}
	return result
}

func reverseSubstitutionIDEA(data []byte, ideaKey []byte) []byte {
	if len(ideaKey) != KeySizeIDEA {
		panic("❌ Ключ должен быть 16 байт")
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeIDEA {
		block := data[i : i+BlockSizeIDEA]
		decrypted := IDEA_DecryptBlock(block, ideaKey)
		copy(result[i:], decrypted)
	}
	return result
}
