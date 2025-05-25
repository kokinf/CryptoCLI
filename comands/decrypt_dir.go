package comands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"CruptoCLI/modules"
)

// DecryptDirCommand — рекурсивно расшифровывает все .enc файлы в каталоге
func DecryptDirCommand(root, output string, masterKey []byte) error {
	fmt.Printf("🔓 Расшифровка каталога: %s → %s\n", root, output)

	err := os.MkdirAll(output, 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать выходной каталог: %v", err)
	}

	filesProcessed := 0
	totalSize := int64(0)

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // не обрабатываем директории
		}

		if !strings.HasSuffix(path, ".enc") {
			return nil // только .enc файлы
		}

		cipherData, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("не удалось прочитать файл %s: %v", path, err)
		}

		plainData, err := modules.MultiLayerDecrypt(cipherData, masterKey)
		if err != nil {
			return fmt.Errorf("ошибка расшифровки файла %s: %v", path, err)
		}

		// Убираем магический заголовок
		plainData = RemoveMagic(plainData)

		rel, _ := filepath.Rel(root, path)
		outputPath := filepath.Join(output, strings.TrimSuffix(rel, ".enc"))

		err = os.MkdirAll(filepath.Dir(outputPath), 0755)
		if err != nil {
			return fmt.Errorf("не удалось создать структуру каталогов: %v", err)
		}

		err = ioutil.WriteFile(outputPath, plainData, 0644)
		if err != nil {
			return fmt.Errorf("не удалось записать файл %s: %v", outputPath, err)
		}

		filesProcessed++
		totalSize += int64(len(plainData))

		fmt.Printf("📄 Расшифрован: %s → %s\n", path, outputPath)

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("✅ Каталог расшифрован: %s → %s\n", root, output)
	fmt.Printf("📁 Обработано файлов: %d | Общий размер: %d байт\n", filesProcessed, totalSize)

	return nil
}
