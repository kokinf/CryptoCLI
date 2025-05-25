package modules

import (
	"bytes"
	"fmt"
)

var (
	mul9  = [256]byte{}
	mul11 = [256]byte{}
	mul13 = [256]byte{}
	mul14 = [256]byte{}
)

// pkcs7Padding — добавляет PKCS#7 padding
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad — удаляет PKCS#7 padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("❌ неверный размер блока")
	}
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("❌ длина данных не кратна размеру блока")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, fmt.Errorf("❌ неверное значение padding'а")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("❌ данные повреждены или padding неверен")
		}
	}
	return data[:len(data)-padding], nil
}

// RemoveMagic — убирает магический заголовок после расшифровки
func RemoveMagic(data []byte) ([]byte, error) {
	if !bytes.HasPrefix(data, []byte(MagicHeader)) {
		return data, nil
	}
	return data[len(MagicHeader):], nil
}

// ShiftRows — сдвиг строк (AES)
func ShiftRows(state []byte) {
	state[1], state[5], state[9], state[13] = state[5], state[9], state[13], state[1]
	state[2], state[6], state[10], state[14] = state[10], state[14], state[2], state[6]
	state[3], state[7], state[11], state[15] = state[15], state[3], state[7], state[11]
}

// InvShiftRows — обратный сдвиг строк (AES)
func InvShiftRows(state []byte) {
	state[5], state[9], state[13], state[1] = state[1], state[5], state[9], state[13]
	state[10], state[14], state[2], state[6] = state[2], state[6], state[10], state[14]
	state[15], state[3], state[7], state[11] = state[3], state[7], state[11], state[15]
}

// AddRoundKey — XOR с раундовым ключом
func AddRoundKey(state, roundKey []byte) {
	if len(roundKey) == 0 {
		panic("❌ Раундовый ключ пуст")
	}
	for i := range state {
		state[i] ^= roundKey[i%len(roundKey)]
	}
}

// splitBlock — деление на 4 слова по 16 бит (IDEA)
func splitBlock(block []byte) [4]uint16 {
	var words [4]uint16
	for i := 0; i < 4; i++ {
		words[i] = uint16(block[2*i])<<8 | uint16(block[2*i+1])
	}
	return words
}

// combineWords — сборка слов обратно в блок (IDEA)
func combineWords(words [4]uint16) []byte {
	block := make([]byte, 8)
	for i := 0; i < 4; i++ {
		block[2*i] = byte(words[i] >> 8)
		block[2*i+1] = byte(words[i])
	}
	return block
}

func xtime(b byte) byte {
	return ((b << 1) ^ ((b >> 7) * 0x1b))
}

func mul2(b byte) byte {
	return xtime(b)
}

func add(a, b uint16) uint16 {
	return a + b
}
