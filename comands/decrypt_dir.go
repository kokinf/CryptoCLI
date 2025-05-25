package comands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"CruptoCLI/modules"
)

func DecryptDirCommand(root, output string, masterKey []byte, verbose bool) error {
	logVerbose(verbose, "🔓 Расшифровка каталога: %s → %s", root, output)

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
		if !strings.HasSuffix(path, ".enc") {
			return nil
		}

		logVerbose(verbose, "📄 Расшифровывается: %s", path)
		cipherData, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("❌ не удалось прочитать файл %s: %v", path, err)
		}

		plainData, err := modules.MultiLayerDecrypt(cipherData, masterKey)
		if err != nil {
			return fmt.Errorf("❌ ошибка расшифровки файла %s: %v", path, err)
		}

		plainData, err = RemoveMagic(plainData)
		if err != nil {
			return fmt.Errorf("❌ ошибка удаления заголовка: %v", err)
		}

		plainData, err = pkcs7Unpad(plainData, modules.BlockSizeAES)
		if err != nil {
			return fmt.Errorf("❌ ошибка PKCS#7 unpad: %v", err)
		}

		rel, _ := filepath.Rel(root, path)
		outPath := filepath.Join(output, strings.TrimSuffix(rel, ".enc"))

		err = os.MkdirAll(filepath.Dir(outPath), 0755)
		if err != nil {
			return fmt.Errorf("❌ не удалось создать структуру каталогов: %v", err)
		}

		err = os.WriteFile(outPath, plainData, 0644)
		if err != nil {
			return fmt.Errorf("❌ не удалось записать файл %s: %v", outPath, err)
		}

		filesProcessed++
		totalSize += int64(len(plainData))
		logVerbose(verbose, "✅ Расшифрован: %s → %s", path, outPath)
		return nil
	})

	if err != nil {
		return err
	}

	logVerbose(verbose, "📁 Обработано файлов: %d | Общий размер: %d байт", filesProcessed, totalSize)
	return nil
}
