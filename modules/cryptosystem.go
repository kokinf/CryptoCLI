package modules

import (
	"bytes"
	"fmt"
)

const MagicHeader = "CRYPTO_UTILITY_V2"

func MultiLayerEncrypt(data, masterKey []byte) ([]byte, error) {
	if len(masterKey) < 32 {
		return nil, fmt.Errorf("❌ мастер-ключ должен быть не менее 32 байт")
	}

	aesKey := masterKey[:16]   // 16 байт
	nushKey := masterKey[4:12] // 8 байт из середины
	kuzKey := masterKey        // 32 байт
	ideaKey := masterKey[8:24] // 16 байт

	result := append([]byte(MagicHeader), data...)

	// Шифруем через каждый слой
	result = permutationAES(result, aesKey)       // P-блок AES
	result = substitutionNUSH(result, nushKey)    // S-блок NUSH
	result = permutationKuznechik(result, kuzKey) // P-блок Kuznechik
	result = substitutionIDEA(result, ideaKey)    // S-блок IDEA

	// Добавляем хэш в начало (аналог HMAC)
	hash := generateHash(result, masterKey)
	result = append(hash[:16], result...)

	return result, nil
}

func MultiLayerDecrypt(data, masterKey []byte) ([]byte, error) {
	if len(masterKey) < 32 {
		return nil, fmt.Errorf("❌ мастер-ключ должен быть не менее 32 байт")
	}

	ideaKey := masterKey[8:24]
	kuzKey := masterKey
	nushKey := masterKey[4:12]
	aesKey := masterKey[:16]

	result := data

	// Проверяем хэш (аналог HMAC)
	if len(result) < 16 {
		return nil, fmt.Errorf("❌ данные слишком маленькие для проверки хэша")
	}
	storedHash := result[:16]
	payload := result[16:]
	expectedHash := generateHash(payload, masterKey)
	if !bytes.Equal(storedHash, expectedHash[:16]) {
		return nil, fmt.Errorf("❌ данные повреждены или хэш неверен")
	}

	result = reverseSubstitutionIDEA(payload, ideaKey)
	result = reversePermutationKuznechik(result, kuzKey)
	result = reverseSubstitutionNUSH(result, nushKey)
	result = reversePermutationAES(result, aesKey)

	// Убираем MagicHeader
	result, err := RemoveMagic(result)
	if err != nil {
		return nil, fmt.Errorf("❌ ошибка удаления заголовка: %v", err)
	}

	// Убираем PKCS#7 padding
	result, err = pkcs7Unpad(result, BlockSizeAES)
	if err != nil {
		return nil, fmt.Errorf("❌ ошибка PKCS#7 unpad: %v", err)
	}

	return result, nil
}

// generateHash — простая функция генерации хэша (аналог HMAC)
func generateHash(data, key []byte) []byte {
	hash := make([]byte, 16)
	for i := 0; i < len(data); i++ {
		hash[i%16] ^= data[i] ^ key[i%len(key)]
	}
	return hash
}
