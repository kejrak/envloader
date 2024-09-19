package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// CheckEncryptedFile checks, if the file is already enrypted.
func CheckEncryptedFile(file string) (bool, error) {
	plainText, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer plainText.Close()

	scanner := bufio.NewScanner(plainText)

	if scanner.Scan() {
		firstLine := scanner.Text()
		return strings.HasPrefix(firstLine, "!envloader | AES-256"), nil
	}

	return false, nil
}

// WriteFile writes byte content to the provided file.
func WriteFile(file string, text []byte) error {

	err := os.WriteFile(file, text, 0644)

	if err != nil {
		log.Printf("write file err: %v", err)
		return err
	}

	return nil
}

// ReadFile reads content from provided file.
func ReadFile(file string) ([]byte, error) {

	plainText, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return plainText, err
}

// WrapBytes wraps byte content.
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

// GetEnv gets environment variable value based on the provided key.
// It cat sets up default value.
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
