// main.go
package main

import (
	"CruptoCLI/comands"
	"CruptoCLI/key"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("❌ Неверное количество аргументов")
		fmt.Println("Использование:")
		fmt.Println("  cryptosystem encrypt file <input> [output]")
		fmt.Println("  cryptosystem decrypt file <input> [output]")
		fmt.Println("  cryptosystem encrypt dir <input_dir> [output_dir]")
		fmt.Println("  cryptosystem decrypt dir <input_dir> [output_dir]")
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

	password := GetPasswordFromUser()
	masterKey := key.GenerateMasterKey(password)

	switch {
	case cmd == "encrypt" && action == "file":
		comands.EncryptFileCommand(src, dst, masterKey)
	case cmd == "decrypt" && action == "file":
		comands.DecryptFileCommand(src, dst, masterKey)
	case cmd == "encrypt" && action == "dir":
		comands.EncryptDirCommand(src, dst, masterKey)
	case cmd == "decrypt" && action == "dir":
		comands.DecryptDirCommand(src, dst, masterKey)
	default:
		fmt.Printf("❌ Неизвестная команда: %s %s\n", cmd, action)
	}
}

// Получение пароля от пользователя (без скрытия символов)
func GetPasswordFromUser() []byte {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("🔑 Введите пароль: ")
	password, _ := reader.ReadString('\n')
	return []byte(strings.TrimSpace(password))
}
