package comands

import (
	"bytes"
	"path/filepath"
)

const MagicHeader = "CRYPTO_UTILITY_V1"

// Проверка наличия магического заголовка
func CheckMagic(data []byte) bool {
	return bytes.HasPrefix(data, []byte(MagicHeader))
}

// Удаление магического заголовка после расшифровки
func RemoveMagic(data []byte) []byte {
	if CheckMagic(data) {
		return data[len(MagicHeader):]
	}
	return data
}

// Получение безопасного выходного пути
func GetDefaultOutputPath(input string) string {
	ext := filepath.Ext(input)
	if ext == ".enc" {
		input = input[:len(input)-4] // убираем .enc
	} else {
		input += ".dec" // добавляем .dec, если нет .enc
	}
	return input
}
