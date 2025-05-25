package comands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"CruptoCLI/modules"
)

// EncryptDirCommand — рекурсивно шифрует все файлы в каталоге
func EncryptDirCommand(root, output string, key []byte) error {
	start := time.Now()
	fmt.Printf("🔒 Начинаю шифрование каталога '%s' → '%s'\n", root, output)

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
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		outputPath := filepath.Join(output, rel+".enc")

		// Чтение файла
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("не удалось прочитать файл %s: %v", path, err)
		}

		// Шифрование
		encrypted := modules.MultiLayerEncrypt(data, key)

		// Запись зашифрованного файла
		err = os.MkdirAll(filepath.Dir(outputPath), 0755)
		if err != nil {
			return fmt.Errorf("не удалось создать путь %s: %v", outputPath, err)
		}

		err = ioutil.WriteFile(outputPath, encrypted, 0644)
		if err != nil {
			return fmt.Errorf("не удалось записать файл %s: %v", outputPath, err)
		}

		filesProcessed++
		totalSize += int64(len(data))
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("✅ Каталог успешно зашифрован\n")
	fmt.Printf("📁 Зашифровано файлов: %d | Общий размер: %d байт\n", filesProcessed, totalSize)
	fmt.Printf("⏱ Время: %v\n", time.Since(start))

	return nil
}
