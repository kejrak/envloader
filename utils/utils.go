package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func CheckEncryptedFile(file string) (bool, error) {
	plainText, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer plainText.Close()

	scanner := bufio.NewScanner(plainText)

	if scanner.Scan() {
		firstLine := scanner.Text()
		return strings.HasPrefix(firstLine, "!envLoader | AES-256"), nil
	}

	return false, nil
}

func WriteFile(file string, text []byte) error {

	err := os.WriteFile(file, text, 0644)

	if err != nil {
		log.Printf("write file err: %v", err)
		return err
	}

	return nil
}

func ReadFile(file string) ([]byte, error) {

	plainText, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return plainText, err
}

func WrapBytes(bytes []byte) string {
	var result string
	var line []byte

	for i, b := range bytes {
		line = append(line, b)

		if len(line) == 70 || i == len(bytes)-1 {
			result += string(line) + "\n"
			line = nil
		}
	}

	return result
}
