package modules

const BlockSizeDES = 8 // 64 бита
const KeySizeDES = 8   // 56 эффективных бит (из 64)

// Полные S-блоки DES
var sBoxDES = [8][4][16]byte{
	{ // S1
		{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
		{0, 15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 1, 13},
		{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
		{13, 0, 11, 7, 4, 9, 1, 10, 12, 2, 5, 8, 6, 12, 15, 1},
	},
	{ // S2
		{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10},
		{3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5},
		{0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15},
		{13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9},
	},
	{ // S3
		{10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8},
		{13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 1, 15},
		{1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12},
		{7, 2, 12, 4, 9, 1, 0, 14, 13, 11, 5, 3, 8, 13, 6, 1},
	},
	{ // S4
		{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
		{14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6},
		{4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14},
		{11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3},
	},
	{ // S5
		{12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11},
		{10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8},
		{9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6},
		{4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13},
	},
	{ // S6
		{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
		{13, 0, 11, 7, 4, 9, 1, 10, 6, 3, 8, 5, 2, 14, 12, 15},
		{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
		{14, 1, 7, 4, 6, 15, 10, 5, 3, 9, 8, 4, 5, 11, 12, 7},
	},
	{ // S7
		{6, 12, 7, 1, 5, 15, 13, 8, 4, 10, 9, 15, 3, 6, 10, 2},
		{13, 1, 2, 15, 8, 13, 4, 8, 6, 10, 15, 9, 0, 14, 5, 3},
		{6, 12, 9, 7, 3, 4, 1, 14, 2, 0, 5, 9, 7, 4, 11, 13},
		{15, 1, 3, 15, 5, 9, 0, 14, 13, 8, 10, 1, 7, 4, 8, 5},
	},
	{ // S8
		{2, 1, 14, 7, 4, 10, 8, 13, 1, 6, 15, 11, 1, 6, 5, 0},
		{8, 2, 11, 1, 12, 14, 15, 5, 0, 9, 7, 3, 10, 5, 0, 2},
		{0, 12, 7, 5, 14, 9, 8, 2, 13, 4, 6, 11, 10, 5, 3, 8},
		{9, 1, 7, 5, 12, 3, 10, 14, 2, 0, 6, 13, 10, 14, 3, 4},
	},
}

// Перестановка начального ключа PC-1
var pc1 = [56]int{
	57, 49, 41, 33, 25, 17, 9,
	1, 58, 50, 42, 34, 26, 18,
	10, 2, 59, 51, 43, 35, 27,
	19, 11, 3, 60, 52, 44, 36,
	63, 55, 47, 39, 31, 23, 15,
	7, 62, 54, 46, 38, 30, 22,
	14, 6, 61, 53, 45, 37, 29,
	21, 13, 5, 28, 20, 12, 4,
}

// Перестановка PC-2 — выборка 48 бит из C и D
var pc2 = [48]int{
	14, 17, 11, 24, 1, 5, 3, 28,
	15, 6, 21, 10, 23, 19, 12, 4,
	26, 8, 16, 7, 27, 20, 13, 41,
	52, 31, 37, 47, 55, 30, 40, 51,
	45, 16, 33, 48, 44, 49, 39, 56,
	34, 53, 46, 42, 50, 36, 29, 32,
}

// Сдвиги на раунд
var shifts = [16]int{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 1, 1}

// IP — начальная перестановка
var ipTable = [64]byte{
	58, 50, 42, 34, 26, 18, 10, 2,
	60, 52, 44, 36, 28, 20, 12, 4,
	62, 54, 46, 38, 30, 22, 14, 6,
	64, 56, 48, 40, 32, 24, 16, 8,
	57, 49, 41, 33, 25, 17, 9, 1,
	59, 51, 43, 35, 27, 19, 11, 3,
	61, 53, 45, 37, 29, 21, 13, 5,
	63, 55, 47, 39, 31, 23, 15, 7,
}

// FP — финальная перестановка
var fpTable = [64]byte{
	40, 8, 48, 16, 56, 24, 64, 32,
	39, 7, 47, 15, 55, 23, 63, 31,
	38, 6, 46, 14, 54, 22, 62, 30,
	37, 5, 45, 13, 53, 21, 61, 29,
	36, 4, 44, 12, 52, 20, 60, 28,
	35, 3, 43, 11, 51, 19, 59, 27,
	34, 2, 42, 10, 50, 18, 58, 26,
	33, 1, 41, 9, 49, 17, 57, 25,
}

// E-расширение — из 32 бит в 48
var eExpansionTable = [48]byte{
	32, 1, 2, 3, 4, 5,
	4, 5, 6, 7, 8, 9,
	8, 9, 10, 11, 12, 13,
	12, 13, 14, 15, 16, 17,
	16, 17, 18, 19, 20, 21,
	20, 21, 22, 23, 24, 25,
	24, 25, 26, 27, 28, 29,
	28, 29, 30, 31, 32, 1,
}

// P-перестановка после S-блоков
var pPermutationTable = [32]byte{
	16, 7, 20, 21, 29, 12, 28, 17,
	1, 15, 23, 26, 5, 18, 31, 10,
	2, 8, 24, 14, 32, 27, 3, 9,
	19, 13, 30, 6, 22, 11, 4, 25,
}

// Разбиение байта на биты
func splitBits(b byte) [8]byte {
	var bits [8]byte
	for i := 0; i < 8; i++ {
		bits[i] = (b >> (7 - i)) & 1
	}
	return bits
}

// Сборка битов в байт
func combineBits(bits [8]byte) byte {
	var b byte
	for i := 0; i < 8; i++ {
		if bits[i] == 1 {
			b |= 1 << (7 - i)
		}
	}
	return b
}

// Битовая перестановка по таблице
func permute(input []byte, table []byte) []byte {
	bitInput := make([]byte, len(table))
	for i, pos := range table {
		byteIdx := (pos - 1) / 8
		bitPos := (pos - 1) % 8
		if byteIdx < uint8(len(input)) && bitPos < 8 {
			bit := (input[byteIdx] >> (7 - bitPos)) & 1
			bitInput[i] = bit
		}
	}

	result := make([]byte, 8)
	for i := 0; i < len(bitInput); i++ {
		result[i/8] |= bitInput[i] << (7 - (i % 8))
	}
	return result
}

// initialPermutation — применение IP
func initialPermutation(block []byte) []byte {
	return permute(block, ipTable[:])
}

// finalPermutation — применение FP
func finalPermutation(block []byte) []byte {
	return permute(block, fpTable[:])
}

// applyPC1 — разбиение ключа через PC1
func applyPC1(key []byte) (uint32, uint32) {
	cBits := make([]byte, 28)
	dBits := make([]byte, 28)

	for i, pos := range pc1 {
		byteIdx := (pos - 1) / 8
		bitPos := (pos - 1) % 8
		if i < 28 {
			cBits[i] = (key[byteIdx] >> (7 - bitPos)) & 1
		} else {
			dBits[i-28] = (key[byteIdx] >> (7 - bitPos)) & 1
		}
	}

	var c, d uint32
	for _, bit := range cBits {
		c <<= 1
		if bit == 1 {
			c |= 1
		}
	}
	for _, bit := range dBits {
		d <<= 1
		if bit == 1 {
			d |= 1
		}
	}
	return c, d
}

// shiftLeft — циклический сдвиг влево
func shiftLeft(val uint32, n int) uint32 {
	return ((val << uint(n)) | (val >> (28 - uint(n)))) & 0x0FFFFFFF
}

// applyPC2 — выборка 48 бит из C и D
func applyPC2(c, d uint32) []byte {
	var subKey [6]byte
	for i, pos := range pc2 {
		var val uint32
		if pos <= 28 {
			val = c
			pos -= 1
		} else {
			val = d
			pos -= 29
		}
		bit := (val >> (28 - pos)) & 1
		subKey[i/8] |= byte(bit) << (7 - uint(i)%8)
	}
	return subKey[:]
}

// generateDESKeySchedule — генерация подключей
func generateDESKeySchedule(key []byte) [][6]byte {
	c, d := applyPC1(key)
	var keySchedule [][6]byte

	for round := 0; round < 16; round++ {
		c = shiftLeft(c, shifts[round])
		d = shiftLeft(d, shifts[round])
		var subkey [6]byte
		copy(subkey[:], applyPC2(c, d))
		keySchedule = append(keySchedule, subkey)
	}
	return keySchedule
}

// eExpansion — расширение 32 бит до 48
func eExpansion(block []byte) []byte {
	bitBlock := make([]byte, 32)
	for i := 0; i < 4; i++ {
		bits := splitBits(block[i])
		copy(bitBlock[i*8:], bits[:])
	}

	e := make([]byte, 48)
	for i := 0; i < 48; i++ {
		pos := eExpansionTable[i] - 1
		e[i] = bitBlock[pos]
	}
	return e
}

// sBoxSubstitute — S-подстановка
func sBoxSubstitute(data []byte) []byte {
	sResult := make([]byte, 4)

	for i := 0; i < 8; i++ {
		row := data[i*6]<<1 | data[i*6+5]
		col := data[i*6+1]<<3 |
			data[i*6+2]<<2 |
			data[i*6+3]<<1 |
			data[i*6+4]

		value := sBoxDES[i][row][col]
		sResult[i/2] <<= 4
		if i%2 == 0 {
			sResult[i/2] |= value >> 4
		} else {
			sResult[i/2] |= value & 0x0f
		}
	}
	return sResult
}

// pPermutation — P-перестановка после S-блоков
func pPermutation(data []byte) []byte {
	bitData := make([]byte, 32)
	for i := 0; i < 4; i++ {
		bits := splitBits(data[i])
		copy(bitData[i*8:], bits[:])
	}

	pResult := make([]byte, 32)
	for i := 0; i < 32; i++ {
		srcBit := pPermutationTable[i] - 1
		pResult[i] = bitData[srcBit]
	}

	result := make([]byte, 4)
	for i := 0; i < 32; i++ {
		result[i/8] |= pResult[i] << (7 - uint(i)%8)
	}
	return result
}

// feistel — функция F в DES
func feistel(block []byte, key []byte) []byte {
	e := eExpansion(block)
	xored := make([]byte, 6)
	for i := range e {
		xored[i] = e[i] ^ key[i]
	}

	s := sBoxSubstitute(xored)
	p := pPermutation(s)
	return p
}

// DES_EncryptBlock — шифрует один блок через DES
func DES_EncryptBlock(block, key []byte) []byte {
	ip := initialPermutation(block)
	L := ip[:4]
	R := ip[4:]

	keySchedule := generateDESKeySchedule(key)

	for round := 0; round < 16; round++ {
		f := feistel(R, keySchedule[round][:])
		nextR := xor(L, f)
		L, R = R, nextR
	}

	fp := finalPermutation(append(R, L...))
	return fp
}

// DES_DecryptBlock — расшифровывает один блок через DES
func DES_DecryptBlock(block, key []byte) []byte {
	ip := initialPermutation(block)
	L := ip[:4]
	R := ip[4:]

	keySchedule := generateDESKeySchedule(key)

	for round := 15; round >= 0; round-- {
		f := feistel(L, keySchedule[round][:])
		nextL := xor(R, f)
		R = L
		L = nextL
	}

	fp := finalPermutation(append(L, R...))
	return fp
}

// permutationDES — шифрует несколько блоков данных через DES
func permutationDES(data []byte, desKey []byte) []byte {
	if len(desKey) != KeySizeDES {
		panic("Неверная длина ключа DES")
	}

	if len(data)%BlockSizeDES != 0 {
		data = pkcs7Padding(data, BlockSizeDES)
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeDES {
		block := data[i : i+BlockSizeDES]
		encrypted := DES_EncryptBlock(block, desKey)
		copy(result[i:], encrypted)
	}
	return result
}

// reversePermutationDES — расшифровывает несколько блоков через DES
func reversePermutationDES(data []byte, desKey []byte) []byte {
	if len(desKey) != KeySizeDES {
		panic("Неверная длина ключа DES")
	}

	if len(data)%BlockSizeDES != 0 {
		data = pkcs7Padding(data, BlockSizeDES)
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeDES {
		block := data[i : i+BlockSizeDES]
		decrypted := DES_DecryptBlock(block, desKey)
		copy(result[i:], decrypted)
	}
	return result
}

// xor — XOR двух блоков
func xor(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}
