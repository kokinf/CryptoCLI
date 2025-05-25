package key

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const fixedSalt = "fixed_salt_128_bit" // Фиксированная соль

func GenerateMasterKey([]byte) []byte {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите пароль: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	masterKey := PBKDF2([]byte(password), []byte(fixedSalt), 10000, 32)
	fmt.Println("Мастер-ключ сгенерирован")
	return masterKey
}
