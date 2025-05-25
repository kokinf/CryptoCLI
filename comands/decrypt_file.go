package comands

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"CruptoCLI/modules"
)

// DecryptFileCommand — CLI команда для расшифровки одного файла
func DecryptFileCommand(input, output string, masterKey []byte, verbose bool) error {
	logVerbose(verbose, "🔓 Расшифровка файла: %s", input)
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return fmt.Errorf("❌ не удалось прочитать файл: %v", err)
	}

	logVerbose(verbose, "🔐 Дешифруем данные")
	start := time.Now()
	decrypted, err := modules.MultiLayerDecrypt(data, masterKey)
	if err != nil {
		return fmt.Errorf("❌ ошибка расшифровки: %v", err)
	}
	elapsed := time.Since(start)
	speed := float64(len(decrypted)*8) / elapsed.Seconds() / 1_000_000
	logVerbose(verbose, "⏱ Скорость расшифровки: %.2f Мбит/сек", speed)

	logVerbose(verbose, "💾 Сохраняем дешифрованный файл: %s", output)
	err = os.WriteFile(output, decrypted, 0644)
	if err != nil {
		return fmt.Errorf("❌ не удалось записать файл: %v", err)
	}

	logVerbose(verbose, "✅ Файл успешно расшифрован")
	return nil
}
