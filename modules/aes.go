package modules

const BlockSizeAES = 16 // AES блок — 128 бит (16 байт)
const KeySizeAES = 16   // AES-128 → ключ 16 байт
const RoundsAES = 10    // 10 раундов для AES-128

// MixColumns — смешивание столбцов
func MixColumns(state []byte) {
	for c := 0; c < 4; c++ {
		b0 := state[c*4+0]
		b1 := state[c*4+1]
		b2 := state[c*4+2]
		b3 := state[c*4+3]

		state[c*4+0] = mul2(b0) ^ b1 ^ mul2(b1) ^ b2 ^ b3
		state[c*4+1] = b0 ^ mul2(b1) ^ b2 ^ mul2(b2) ^ b3
		state[c*4+2] = b0 ^ b1 ^ mul2(b2) ^ b3 ^ mul2(b3)
		state[c*4+3] = mul2(b0) ^ b1 ^ b2 ^ mul2(b3)
	}
}

// InvMixColumns — обратное смешивание столбцов
func InvMixColumns(state []byte) {
	for c := 0; c < 4; c++ {
		b0 := state[c*4+0]
		b1 := state[c*4+1]
		b2 := state[c*4+2]
		b3 := state[c*4+3]

		state[c*4+0] = mul9[b0] ^ mul11[b1] ^ mul13[b2] ^ mul14[b3]
		state[c*4+1] = mul14[b0] ^ mul9[b1] ^ mul11[b2] ^ mul13[b3]
		state[c*4+2] = mul13[b0] ^ mul14[b1] ^ mul9[b2] ^ mul11[b3]
		state[c*4+3] = mul11[b0] ^ mul13[b1] ^ mul14[b2] ^ mul9[b3]
	}
}

// xtime — умножение на x в GF(2^8)
func xtime(b byte) byte {
	return ((b << 1) ^ ((b >> 7) * 0x1b))
}

// mul2 — умножение на 2 в поле Галуа
func mul2(b byte) byte { return xtime(b) }

// Таблицы для InvMixColumns (AES)
var (
	mul9  = [256]byte{}
	mul11 = [256]byte{}
	mul13 = [256]byte{}
	mul14 = [256]byte{}
)

func init() {
	for i := 0; i < 256; i++ {
		b := byte(i)
		mul9[i] = mul2(mul2(mul2(b)^b)) ^ b
		mul11[i] = mul2(mul2(mul2(b) ^ b ^ mul2(b)))
		mul13[i] = mul2(mul2(mul2(b)^b)) ^ mul2(b) ^ b
		mul14[i] = mul2(mul2(mul2(b)^b)) ^ b
	}
}

// KeyExpansion — расширение ключа
func KeyExpansion(key []byte) [][]byte {
	if len(key) != KeySizeAES {
		panic("Неверная длина ключа")
	}

	keySchedule := make([][]byte, 4*(RoundsAES+1))
	for i := range keySchedule {
		keySchedule[i] = make([]byte, 4)
	}

	// Инициализация первых слов
	for i := 0; i < 4; i++ {
		copy(keySchedule[i], key[i*4:i*4+4])
	}

	rconIndex := 1
	for i := 4; i < len(keySchedule); i++ {
		temp := keySchedule[i-1]
		if i%4 == 0 {
			temp = append(temp[1:], temp[0]) // RotWord
			for j := range temp {
				temp[j] = sBox[temp[j]] // SubWord
			}
			temp[0] ^= rcon[rconIndex]
			rconIndex++
		} else if i > 0 && i%4 == 0 {
			for j := range temp {
				temp[j] = sBox[temp[j]]
			}
		}

		for j := 0; j < 4; j++ {
			keySchedule[i][j] = keySchedule[i-4][j] ^ temp[j]
		}
	}

	return keySchedule
}

var rcon = [255]byte{
	0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80,
	0x1b, 0x36, 0x6c, 0xd8, 0xab, 0x4d, 0x9a, 0x2f,
	0x5e, 0xbc, 0x63, 0xc6, 0x97, 0x35, 0x6a, 0xda,
	0xa3, 0x5d, 0xb4, 0x73, 0xe6, 0xd7, 0xb5, 0x71,
}

// PKCS#7 padding

// AES_EncryptBlock — шифрует один 16-байтовый блок через AES
func AES_EncryptBlock(block, key []byte) []byte {
	state := make([]byte, BlockSizeAES)
	copy(state, block)

	keySchedule := KeyExpansion(key)

	AddRoundKey(state, keySchedule[0])

	for round := 1; round < 10; round++ {
		SubBytes(state)
		ShiftRows(state)
		MixColumns(state)
		AddRoundKey(state, keySchedule[round])
	}

	// Последний раунд без смешивания
	SubBytes(state)
	ShiftRows(state)
	AddRoundKey(state, keySchedule[10])

	return state
}

// AES_DecryptBlock — расшифровывает один 16-байтовый блок
func AES_DecryptBlock(block, key []byte) []byte {
	state := make([]byte, BlockSizeAES)
	copy(state, block)

	keySchedule := KeyExpansion(key)

	AddRoundKey(state, keySchedule[10])

	for round := 9; round >= 0; round-- {
		InvShiftRows(state)
		InvSubBytes(state)
		AddRoundKey(state, keySchedule[round])
		if round > 0 {
			InvMixColumns(state)
		}
	}

	return state
}

// permutationAES — шифрует несколько блоков через AES
func permutationAES(data []byte, aesKey []byte) []byte {
	if len(data)%BlockSizeAES != 0 {
		data = pkcs7Padding(data, BlockSizeAES)
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeAES {
		block := data[i : i+BlockSizeAES]
		encrypted := AES_EncryptBlock(block, aesKey)
		copy(result[i:], encrypted)
	}
	return result
}

// reversePermutationAES — расшифровывает несколько блоков через AES
func reversePermutationAES(data []byte, aesKey []byte) []byte {
	if len(data)%BlockSizeAES != 0 {
		data = pkcs7Padding(data, BlockSizeAES)
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeAES {
		block := data[i : i+BlockSizeAES]
		decrypted := AES_DecryptBlock(block, aesKey)
		copy(result[i:], decrypted)
	}
	return result
}
