package comands

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"CruptoCLI/modules"
)

// EncryptFileCommand — CLI команда для шифрования одного файла
func EncryptFileCommand(input, output string, masterKey []byte, verbose bool) error {
	logVerbose(verbose, "📄 Зачитываем файл: %s", input)
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return fmt.Errorf("❌ не удалось прочитать файл: %v", err)
	}

	logVerbose(verbose, "🔒 Шифруем через все слои")
	start := time.Now()
	encrypted, err := modules.MultiLayerEncrypt(data, masterKey)
	if err != nil {
		return fmt.Errorf("❌ ошибка шифрования: %v", err)
	}
	elapsed := time.Since(start)
	speed := float64(len(data)*8) / elapsed.Seconds() / 1_000_000
	logVerbose(verbose, "⏱ Скорость шифрования: %.2f Мбит/сек", speed)

	logVerbose(verbose, "💾 Сохраняем зашифрованный файл: %s", output)
	err = os.WriteFile(output, encrypted, 0644)
	if err != nil {
		return fmt.Errorf("❌ не удалось записать файл: %v", err)
	}

	logVerbose(verbose, "✅ Файл успешно зашифрован")
	return nil
}
