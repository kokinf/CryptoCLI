package comands

import (
	"bytes"
	"fmt"
	"strings"
)

const MagicHeader = "CRYPTO_UTILITY_V2"

// GetDefaultOutputPath — безопасное имя выходного файла
func GetDefaultOutputPath(input string) string {
	if strings.HasSuffix(input, ".enc") {
		return input[:len(input)-4]
	}
	return input + ".enc"
}

// RemoveMagic — удаляет магический заголовок после расшифровки
func RemoveMagic(data []byte) ([]byte, error) {
	if len(data) < len(MagicHeader) {
		return data, nil
	}
	if !bytes.HasPrefix(data, []byte(MagicHeader)) {
		return data, nil
	}
	return data[len(MagicHeader):], nil
}

// pkcs7Unpad — убирает padding после расшифровки
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("❌ неверный размер блока")
	}
	if len(data)%blockSize != 0 {
		return nil, fmt.Errorf("❌ длина данных не кратна размеру блока")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, fmt.Errorf("❌ неверное значение padding'а")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("❌ неверный padding")
		}
	}
	return data[:len(data)-padding], nil
}

// zeroBytes — обнуляет байты для безопасности
func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// logVerbose — вывод только если включен режим verbose
func logVerbose(verbose bool, format string, a ...interface{}) {
	if verbose {
		fmt.Printf("[VERBOSE] "+format+"\n", a...)
	}
}
