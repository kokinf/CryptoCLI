// modules/nush.go
package modules

const BlockSizeNUSH = 16 // NUSH работает с блоками по 128 бит (16 байт)
const KeySizeNUSH = 8    // Ключ NUSH — 64 бита (8 байт)
const RoundsNUSH = 12    // 12 раундов для NUSH

// NUSH_EncryptBlock — шифрует один блок NUSH
func NUSH_EncryptBlock(block, key []byte) []byte {
	if len(block) != BlockSizeNUSH {
		panic("Неверная длина блока NUSH")
	}
	if len(key) != KeySizeNUSH {
		panic("Неверная длина ключа NUSH")
	}

	state := make([]byte, BlockSizeNUSH)
	copy(state, block)

	keySchedule := generateNUSHKeySchedule(key)

	// Основные раунды
	for round := 0; round < RoundsNUSH; round++ {
		SubBytes(state)
		ShiftRows(state)
		AddRoundKey(state, keySchedule[round])
	}

	return state
}

// NUSH_DecryptBlock — расшифровывает один блок NUSH
func NUSH_DecryptBlock(block, key []byte) []byte {
	if len(block) != BlockSizeNUSH {
		panic("Неверная длина блока NUSH")
	}
	if len(key) != KeySizeNUSH {
		panic("Неверная длина ключа NUSH")
	}

	state := make([]byte, BlockSizeNUSH)
	copy(state, block)

	keySchedule := generateNUSHKeySchedule(key)

	// Обратные раунды
	for round := RoundsNUSH - 1; round >= 0; round-- {
		AddRoundKey(state, keySchedule[round])
		InvShiftRows(state)
		InvSubBytes(state)
	}

	return state
}

// substitutionNUSH — шифрует несколько блоков через NUSH
func substitutionNUSH(data []byte, nushKey []byte) []byte {
	if len(nushKey) != KeySizeNUSH {
		panic("Неверная длина ключа NUSH")
	}

	if len(data)%BlockSizeNUSH != 0 {
		data = pkcs7Padding(data, BlockSizeNUSH)
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeNUSH {
		block := data[i : i+BlockSizeNUSH]
		encrypted := NUSH_EncryptBlock(block, nushKey)
		copy(result[i:], encrypted)
	}
	return result
}

// reverseSubstitutionNUSH — расшифровывает несколько блоков через NUSH
func reverseSubstitutionNUSH(data []byte, nushKey []byte) []byte {
	if len(nushKey) != KeySizeNUSH {
		panic("Неверная длина ключа NUSH")
	}

	if len(data)%BlockSizeNUSH != 0 {
		data = pkcs7Padding(data, BlockSizeNUSH)
	}

	result := make([]byte, len(data))
	for i := 0; i < len(data); i += BlockSizeNUSH {
		block := data[i : i+BlockSizeNUSH]
		decrypted := NUSH_DecryptBlock(block, nushKey)
		copy(result[i:], decrypted)
	}
	return result
}

// Генерация раундовых ключей для NUSH (упрощённая версия)
func generateNUSHKeySchedule(key []byte) [][]byte {
	keySchedule := make([][]byte, RoundsNUSH)
	for r := 0; r < RoundsNUSH; r++ {
		start := r % len(key)
		end := start + KeySizeNUSH
		k := key[start:end]
		if end > len(key) {
			k = append(k, key[:end-len(key)]...)
		}
		keySchedule[r] = k
	}
	return keySchedule
}
