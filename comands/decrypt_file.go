package comands

import (
	"fmt"
	"io/ioutil"
	"os"

	"CruptoCLI/modules"
)

// DecryptFileCommand — расшифровывает один файл
func DecryptFileCommand(input, output string, key []byte) error {
	fmt.Printf("🔓 Расшифровка файла: %s → %s\n", input, output)

	cipherData, err := ioutil.ReadFile(input)
	if err != nil {
		return fmt.Errorf("не удалось прочитать зашифрованный файл: %v", err)
	}

	plainData, err := modules.MultiLayerDecrypt(cipherData, key)
	if err != nil {
		return fmt.Errorf("ошибка расшифровки: %v", err)
	}

	// Убираем MagicHeader
	plainData = RemoveMagic(plainData)

	err = os.WriteFile(output, plainData, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать расшифрованный файл: %v", err)
	}

	fmt.Printf("✅ Файл успешно расшифрован: %s → %s\n", input, output)
	return nil
}
