package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"CruptoCLI/comands"
	"CruptoCLI/key"
)

func getPasswordFromUser() []byte {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("🔑 Введите пароль: ")
	password, _ := reader.ReadString('\n')
	return []byte(strings.TrimSpace(password))
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "\n🔐 CruptoCLI — многослойная криптосистема\n")
	fmt.Fprintf(os.Stderr, "Использование:\n")
	fmt.Fprintf(os.Stderr, "  cryptosystem encrypt file <input> [output] [-v]\n")
	fmt.Fprintf(os.Stderr, "  cryptosystem decrypt file <input> [output] [-v]\n")
	fmt.Fprintf(os.Stderr, "  cryptosystem encrypt dir <input_dir> [output_dir] [-v]\n")
	fmt.Fprintf(os.Stderr, "  cryptosystem decrypt dir <input_dir> [output_dir] [-v]\n")
	fmt.Fprintf(os.Stderr, "\nПримеры:\n")
	fmt.Fprintf(os.Stderr, "  Зашифровать файл:\n    ./cryptosystem encrypt file secret.txt secret.enc -v\n")
	fmt.Fprintf(os.Stderr, "  Расшифровать файл:\n    ./cryptosystem decrypt file secret.enc decrypted.txt -v\n")
	fmt.Fprintf(os.Stderr, "  Зашифровать каталог:\n    ./cryptosystem encrypt dir data encrypted_data -v\n")
	fmt.Fprintf(os.Stderr, "  Расшифровать каталог:\n    ./cryptosystem decrypt dir encrypted_data decrypted_data -v\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 4 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		showUsage()
		return
	}

	cmd := os.Args[1]
	action := os.Args[2]
	src := os.Args[3]

	var dst string
	if len(os.Args) > 4 {
		dst = os.Args[4]
	} else {
		dst = comands.GetDefaultOutputPath(src)
	}

	password := getPasswordFromUser()

	masterKey := key.GenerateMasterKey(password)

	if len(masterKey) < 32 {
		fmt.Fprintf(os.Stderr, "❌ Мастер-ключ слишком короткий\n")
		os.Exit(1)
	}

	// Выбор действия
	switch {
	case cmd == "encrypt" && action == "file":
		err := comands.EncryptFileCommand(src, dst, masterKey, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка шифрования файла: %v\n", err)
			os.Exit(1)
		}
	case cmd == "decrypt" && action == "file":
		err := comands.DecryptFileCommand(src, dst, masterKey, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка расшифровки файла: %v\n", err)
			os.Exit(1)
		}
	case cmd == "encrypt" && action == "dir":
		err := comands.EncryptDirCommand(src, dst, masterKey, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка шифрования каталога: %v\n", err)
			os.Exit(1)
		}
	case cmd == "decrypt" && action == "dir":
		err := comands.DecryptDirCommand(src, dst, masterKey, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка расшифровки каталога: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "❌ Неверная команда: %s %s\n", cmd, action)
		showUsage()
	}
}
