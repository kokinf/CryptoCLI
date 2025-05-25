package comands

import (
	"CruptoCLI/modules"
	"fmt"
	"io/ioutil"
	"time"
)

func EncryptFileCommand(input, output string, key []byte) error {
	start := time.Now()

	data, err := ioutil.ReadFile(input)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	encrypted := modules.MultiLayerEncrypt(data, key)

	err = ioutil.WriteFile(output, encrypted, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать файл: %v", err)
	}

	fmt.Printf("⏱ Шифровка завершена за %v\n", time.Since(start))
	return nil
}
