package modules

const BlockSizeKuznechik = 16
const KeySizeKuznechik = 32
const RoundsKuznechik = 9

var sBoxKuz = [256]byte{
	0xFC, 0xEE, 0xDD, 0x11, 0xCF, 0x6B, 0x37, 0xC8,
	0xB4, 0x28, 0xAD, 0x0D, 0xF8, 0xEB, 0x2F, 0xD4,
	0x19, 0xFF, 0x10, 0x74, 0x4A, 0x92, 0x97, 0x5E,
	0xFB, 0x8C, 0x26, 0xB3, 0xEA, 0x20, 0x35, 0xDF,
	0x44, 0xA4, 0xE5, 0xC0, 0x85, 0x32, 0x3B, 0x94,
	0x62, 0x00, 0xBB, 0x71, 0x24, 0x50, 0x05, 0x14,
	0xA7, 0x55, 0x3E, 0x15, 0x2D, 0xE1, 0x21, 0x95,
	0x45, 0xA5, 0x64, 0xAF, 0xF2, 0xE7, 0x99, 0xFA,
	0x72, 0x13, 0xB5, 0xFE, 0x8B, 0x96, 0xD2, 0xBF,
	0xBC, 0x6C, 0x01, 0x8D, 0x12, 0x40, 0x48, 0x1E,
	0xDE, 0x5A, 0x04, 0x77, 0x3B, 0x10, 0xBA, 0x7C,
	0x6F, 0x4F, 0x07, 0x0A, 0xB7, 0x84, 0xC7, 0x23,
	0xC3, 0x18, 0x96, 0x05, 0x9A, 0x07, 0x12, 0x80,
	0xE2, 0xEB, 0x27, 0xB2, 0x75, 0x09, 0x83, 0x2C,
	0x1A, 0x1B, 0x6E, 0x5A, 0xA0, 0x52, 0x3B, 0xD6,
	0xB3, 0x29, 0xE3, 0x2F, 0x84, 0x53, 0xD1, 0x00,
	0xED, 0x20, 0xFC, 0xB1, 0x5B, 0x6A, 0xCB, 0xBE,
	0x39, 0x4A, 0x4C, 0x58, 0xCF, 0xD0, 0xEF, 0xAA,
	0xFB, 0x43, 0x4D, 0x33, 0x85, 0x45, 0xF9, 0x02,
	0x7F, 0x50, 0x3C, 0x9F, 0xA8, 0x51, 0xA3, 0x40,
	0x8F, 0x92, 0x9D, 0x38, 0xF5, 0xBC, 0xB6, 0xDA,
	0x21, 0x10, 0xFF, 0xF3, 0xD2, 0xCD, 0x0C, 0x13,
	0xEC, 0x5F, 0x97, 0x44, 0x17, 0xC4, 0xA7, 0x7E,
	0x3D, 0x64, 0x5D, 0x19, 0x73, 0x60, 0x81, 0x4F,
	0xDC, 0x22, 0x2A, 0x90, 0x88, 0x46, 0xEE, 0xB8,
	0x14, 0xDE, 0x5E, 0x0B, 0xDB, 0xE0, 0x32, 0x3A,
	0x0A, 0x49, 0x06, 0x24, 0x5C, 0xC2, 0xD3, 0xAC,
	0x62, 0x91, 0x95, 0xE4, 0x79, 0xE7, 0xC8, 0x37,
	0x6D, 0x8D, 0xD5, 0x4E, 0xA9, 0x6C, 0x56, 0xF4,
	0xEA, 0x65, 0x7A, 0xAE, 0x08, 0xBA, 0x78, 0x25,
	0x2E, 0x1C, 0xA6, 0xB4, 0xC6, 0xE8, 0xDD, 0x74,
	0x1F, 0x4B, 0xBD, 0x8B, 0x8A, 0x70, 0x3E, 0xB5,
}

func SubBytesKuz(state []byte) {
	for i := range state {
		state[i] = sBoxKuz[state[i]]
	}
}

// InvSubBytesKuz — обратная замена байтов
func InvSubBytesKuz(state []byte) {
	for i := range state {
		state[i] = sBoxKuz[state[i]]
	}
}

func ShiftRowsKuz(state []byte) {
	state[1], state[5], state[9], state[13] = state[5], state[9], state[13], state[1]
	state[2], state[6], state[10], state[14] = state[10], state[14], state[2], state[6]
	state[3], state[7], state[11], state[15] = state[15], state[3], state[7], state[11]
}

// InvShiftRowsKuz — обратный сдвиг строк
func InvShiftRowsKuz(state []byte) {
	state[5], state[9], state[13], state[1] = state[1], state[5], state[9], state[13]
	state[10], state[14], state[2], state[6] = state[2], state[6], state[10], state[14]
	state[15], state[3], state[7], state[11] = state[3], state[7], state[11], state[15]
}

// MixColumnsKuz — смешивание столбцов
func MixColumnsKuz(data []byte) []byte {
	state := make([]byte, 16)
	copy(state, data)

	mds := [4][4]byte{
		{0x02, 0x03, 0x01, 0x01},
		{0x01, 0x02, 0x03, 0x01},
		{0x01, 0x01, 0x02, 0x03},
		{0x03, 0x01, 0x01, 0x02},
	}

	for c := 0; c < 4; c++ {
		b0 := state[c*4+0]
		b1 := state[c*4+1]
		b2 := state[c*4+2]
		b3 := state[c*4+3]

		state[c*4+0] = gF128Mult(mds[0][0], b0) ^
			gF128Mult(mds[0][1], b1) ^
			gF128Mult(mds[0][2], b2) ^
			gF128Mult(mds[0][3], b3)

		state[c*4+1] = gF128Mult(mds[1][0], b0) ^
			gF128Mult(mds[1][1], b1) ^
			gF128Mult(mds[1][2], b2) ^
			gF128Mult(mds[1][3], b3)

		state[c*4+2] = gF128Mult(mds[2][0], b0) ^
			gF128Mult(mds[2][1], b1) ^
			gF128Mult(mds[2][2], b2) ^
			gF128Mult(mds[2][3], b3)

		state[c*4+3] = gF128Mult(mds[3][0], b0) ^
			gF128Mult(mds[3][1], b1) ^
			gF128Mult(mds[3][2], b2) ^
			gF128Mult(mds[3][3], b3)
	}

	return state
}

// InvMixColumnsKuz — обратное смешивание столбцов
func InvMixColumnsKuz(data []byte) []byte {
	state := make([]byte, 16)
	copy(state, data)

	// Обратная матрица MDS
	invMDS := [4][4]byte{
		{0x0E, 0x0B, 0x0D, 0x09},
		{0x09, 0x0E, 0x0B, 0x0D},
		{0x0D, 0x09, 0x0E, 0x0B},
		{0x0B, 0x0D, 0x09, 0x0E},
	}

	for c := 0; c < 4; c++ {
		b0 := state[c*4+0]
		b1 := state[c*4+1]
		b2 := state[c*4+2]
		b3 := state[c*4+3]

		state[c*4+0] = gF128Mult(invMDS[0][0], b0) ^
			gF128Mult(invMDS[0][1], b1) ^
			gF128Mult(invMDS[0][2], b2) ^
			gF128Mult(invMDS[0][3], b3)

		state[c*4+1] = gF128Mult(invMDS[1][0], b0) ^
			gF128Mult(invMDS[1][1], b1) ^
			gF128Mult(invMDS[1][2], b2) ^
			gF128Mult(invMDS[1][3], b3)

		state[c*4+2] = gF128Mult(invMDS[2][0], b0) ^
			gF128Mult(invMDS[2][1], b1) ^
			gF128Mult(invMDS[2][2], b2) ^
			gF128Mult(invMDS[2][3], b3)

		state[c*4+3] = gF128Mult(invMDS[3][0], b0) ^
			gF128Mult(invMDS[3][1], b1) ^
			gF128Mult(invMDS[3][2], b2) ^
			gF128Mult(invMDS[3][3], b3)
	}
	return state
}

func gF128Mult(a, b byte) byte {
	var p byte
	for i := 0; i < 8; i++ {
		if b&0x80 != 0 {
			p ^= a
		}
		highBit := a >> 7
		a <<= 1
		if highBit == 1 {
			a ^= 0x1b
		}
		b <<= 1
	}
	return p
}

func AddRoundKeyKuz(state, roundKey []byte) {
	if len(roundKey) < 16 {
		panic("❌ Раундовый ключ слишком мал")
	}
	for i := 0; i < 16; i++ {
		state[i] ^= roundKey[i]
	}
}

// generateKuzKeySchedule — расширяет мастер-ключ в раундовые ключи
func generateKuzKeySchedule(key []byte) [][16]byte {
	if len(key) != KeySizeKuznechik {
		panic("❌ Неверная длина мастер-ключа")
	}

	const totalKeys = RoundsKuznechik + 1 // 10 раундовых ключей
	keySchedule := make([][16]byte, totalKeys)

	copy(keySchedule[0][:], key[0:16])
	copy(keySchedule[1][:], key[16:32])

	c := [16]byte{}
	for i := 2; i < totalKeys; i++ {
		for j := 0; j < 16; j++ {
			keySchedule[i][j] = keySchedule[i-1][j] ^ c[j]
		}
	}
	return keySchedule
}

func Kuznechik_EncryptBlock(block, key []byte) []byte {
	if len(block) != BlockSizeKuznechik || len(key) != KeySizeKuznechik {
		panic("❌ Неверная длина блока или ключа")
	}

	state := make([]byte, BlockSizeKuznechik)
	copy(state, block)

	keySchedule := generateKuzKeySchedule(key)

	AddRoundKeyKuz(state, keySchedule[0][:])
	for round := 1; round < RoundsKuznechik; round++ {
		SubBytesKuz(state)
		ShiftRowsKuz(state)
		MixColumnsKuz(state)
		AddRoundKeyKuz(state, keySchedule[round][:])
	}
	SubBytesKuz(state)
	ShiftRowsKuz(state)
	AddRoundKeyKuz(state, keySchedule[RoundsKuznechik][:])

	return state
}

func Kuznechik_DecryptBlock(block, key []byte) []byte {
	if len(block) != BlockSizeKuznechik || len(key) != KeySizeKuznechik {
		panic("❌ Неверная длина блока или ключа")
	}

	state := make([]byte, BlockSizeKuznechik)
	copy(state, block)

	keySchedule := generateKuzKeySchedule(key)

	AddRoundKeyKuz(state, keySchedule[RoundsKuznechik][:])
	for round := RoundsKuznechik - 1; round >= 1; round-- {
		InvShiftRowsKuz(state)
		InvSubBytesKuz(state)
		AddRoundKeyKuz(state, keySchedule[round][:])
		InvMixColumnsKuz(state)
	}
	InvShiftRowsKuz(state)
	InvSubBytesKuz(state)
	AddRoundKeyKuz(state, keySchedule[0][:])

	return state
}

// permutationKuznechik — шифрует несколько блоков
func permutationKuznechik(data []byte, key []byte) []byte {
	if len(key) != KeySizeKuznechik {
		panic("❌ Неверная длина ключа")
	}
	data = pkcs7Padding(data, BlockSizeKuznechik)
	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeKuznechik {
		block := data[i : i+BlockSizeKuznechik]
		encrypted := Kuznechik_EncryptBlock(block, key)
		copy(result[i:], encrypted)
	}
	return result
}

// reversePermutationKuznechik — расшифровывает несколько блоков
func reversePermutationKuznechik(data []byte, key []byte) []byte {
	if len(key) != KeySizeKuznechik {
		panic("❌ Неверная длина ключа")
	}
	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeKuznechik {
		block := data[i : i+BlockSizeKuznechik]
		decrypted := Kuznechik_DecryptBlock(block, key)
		copy(result[i:], decrypted)
	}
	return result
}
