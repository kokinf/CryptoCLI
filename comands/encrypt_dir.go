package comands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"CruptoCLI/modules"
)

func EncryptDirCommand(root, output string, masterKey []byte, verbose bool) error {
	logVerbose(verbose, "📁 Начинаю шифрование каталога: %s → %s", root, output)

	err := os.MkdirAll(output, 0755)
	if err != nil {
		return fmt.Errorf("❌ не удалось создать выходной каталог: %v", err)
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
		outPath := filepath.Join(output, rel+".enc")

		logVerbose(verbose, "📄 Зашифровывается: %s", path)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("❌ не удалось прочитать файл %s: %v", path, err)
		}

		encrypted, err := modules.MultiLayerEncrypt(data, masterKey)
		if err != nil {
			return fmt.Errorf("❌ ошибка шифрования файла %s: %v", path, err)
		}

		err = os.MkdirAll(filepath.Dir(outPath), 0755)
		if err != nil {
			return fmt.Errorf("❌ не удалось создать структуру каталогов: %v", err)
		}

		err = os.WriteFile(outPath, encrypted, 0644)
		if err != nil {
			return fmt.Errorf("❌ не удалось записать файл %s: %v", outPath, err)
		}

		filesProcessed++
		totalSize += int64(len(data))
		logVerbose(verbose, "✅ Зашифрован: %s → %s", path, outPath)
		return nil
	})

	if err != nil {
		return err
	}

	logVerbose(verbose, "📁 Обработано файлов: %d | Общий размер: %d байт", filesProcessed, totalSize)
	return nil
}
