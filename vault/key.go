package vault

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kejrak/envLoader/utils"
	"golang.org/x/term"
)

const (
	minPasswordLength = 4
	keyLength         = 32
)

// KeyManager is an interface for managing encryption keys.
type KeyManager interface {
	getKey() (string, error)
}

// KeyType represents the type of encryption key and its source.
type KeyType struct {
	keyFile            string
	keyString          string
	encryptionRequired bool
}

// getKey gets the encryption key based on the KeyType configuration.
func (km *KeyType) getKey() (string, error) {

	if km.keyString != "" {
		return getKeyFromString(km.keyString)
	}

	if km.keyFile != "" {
		return getKeyFromFile(km.keyFile)
	}

	return getKeyFromPrompt(km.encryptionRequired)
}

// getKeyFromString gets the encryption key from a string.
func getKeyFromString(keyString string) (string, error) {
	key, err := fillKeyString(keyString, keyLength)
	if err != nil {
		return "", err
	}
	return key, nil
}

// getKeyFromFile gets the encryption key from provided file.
func getKeyFromFile(keyFile string) (string, error) {

	password, err := utils.ReadFile(keyFile)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	key, err := fillKeyString(string(password), keyLength)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

// getKeyFromPrompt gets the encryption key from user input.
func getKeyFromPrompt(encryptionRequired bool) (string, error) {
	if encryptionRequired {
		bytePassword, err := readPassword("New password: ")
		if err != nil {
			return "", fmt.Errorf("\nfailed to get key from prompt: %w", err)
		}

		checkPassword, err := readPassword("\nRepeat password: ")
		fmt.Println()
		if err != nil {
			return "", fmt.Errorf("failed to get key from prompt: %w", err)
		}

		if bytePassword != checkPassword {
			return "", errors.New("passwords don't match")
		}

		password := string(bytePassword)
		key, err := fillKeyString(strings.TrimSpace(password), keyLength)
		if err != nil {
			return "", err
		}

		return key, nil
	}
	password, err := readPassword("Password: ")
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to get key from prompt: %w", err)
	}
	key, err := fillKeyString(password, keyLength)
	if err != nil {
		return "", fmt.Errorf("failed to get key from prompt: %w", err)
	}
	return key, nil
}

// readPassword reads the user input from STDIN.
func readPassword(promt string) (string, error) {
	fmt.Print(promt)
	stdin := int(syscall.Stdin)
	oldState, err := term.GetState(stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	defer term.Restore(stdin, oldState)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		for _ = range sigch {
			term.Restore(stdin, oldState)
			fmt.Print("\nReceived interrupt signal.\n")
			os.Exit(1)
		}
	}()

	password, err := term.ReadPassword(stdin)
	if err != nil {
		return "", err
	}
	return string(password), nil
}

// fillKeyString fills up the key to provided maximum keyLength.
func fillKeyString(key string, keyLength int) (string, error) {
	if len(key) <= minPasswordLength {
		return "", errors.New("provided password is too short")
	}

	bytes := []byte(key)
	iter := len(bytes)

	for i := 0; i < keyLength-iter; i++ {
		bytes = append(bytes, byte(0))
	}

	password := string(bytes)
	return password, nil
}
