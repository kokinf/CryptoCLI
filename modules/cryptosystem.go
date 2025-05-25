package modules

import (
	"bytes"
	"fmt"
)

const MagicHeader = "CRYPTO_UTILITY_V1"

// MultiLayerEncrypt — многослойное шифрование: AES → NUSH → DES → IDEA
func MultiLayerEncrypt(data, masterKey []byte) []byte {
	aesKey := masterKey[:16]   // 128 бит для AES
	nushKey := masterKey[4:12] // 8 байт = 64 бита для NUSH
	desKey := masterKey[:8]    // 8 байт = 64 бита для DES
	ideaKey := masterKey[8:24] // 16 байт = 128 бит для IDEA

	result := data

	// Добавляем MagicHeader для проверки пароля при расшифровке
	result = append([]byte(MagicHeader), result...)

	// Шифруем через каждый слой
	result = permutationAES(result, aesKey)    // AES (перестановка)
	result = substitutionNUSH(result, nushKey) // NUSH (подстановка)
	result = permutationDES(result, desKey)    // DES (перестановка)
	result = substitutionIDEA(result, ideaKey) // IDEA (подстановка)

	return result
}

// MultiLayerDecrypt — обратная дешифровка: IDEA ← DES ← NUSH ← AES
func MultiLayerDecrypt(data, masterKey []byte) ([]byte, error) {
	ideaKey := masterKey[8:24]
	desKey := masterKey[:8]
	nushKey := masterKey[4:12]
	aesKey := masterKey[:16]

	result := data

	// Дешифровка в обратном порядке
	result = reverseSubstitutionIDEA(result, ideaKey) // IDEA⁻¹
	result = reversePermutationDES(result, desKey)    // DES⁻¹
	result = reverseSubstitutionNUSH(result, nushKey) // NUSH⁻¹
	result = reversePermutationAES(result, aesKey)    // AES⁻¹

	// Проверяем MagicHeader
	if len(result) < len(MagicHeader) {
		return nil, fmt.Errorf("слишком мало данных для проверки заголовка")
	}
	if !bytes.HasPrefix(result, []byte(MagicHeader)) {
		return nil, fmt.Errorf("неверный пароль или повреждённые данные")
	}

	// Убираем MagicHeader и возвращаем оригинальные данные
	return result[len(MagicHeader):], nil
}
