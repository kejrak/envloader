package vault

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kejrak/envLoader/utils"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/term"
)

// KeyManager is an interface for managing encryption keys.
type KeyManager interface {
	getKey() (string, error)
	deriveKey(key, salt []byte) ([]byte, []byte, error)
}

// KeyType represents the type of encryption key and its source.
type KeyType struct {
	keyFile            string
	keyString          string
	encryptionRequired bool
}

// getKey gets the encryption key based on the KeyType configuration.
func (km *KeyType) getKey() ([]byte, error) {
	if km.keyString != "" {
		return getKeyFromString(km.keyString)
	}

	if km.keyFile != "" {
		return getKeyFromFile(km.keyFile)
	}

	return getKeyFromPrompt(km.encryptionRequired)
}

// deriveKey gets the encryption key based on salt configurations.
func (km *KeyType) deriveKey(key, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(key, salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

// getKeyFromString gets the encryption key from a string.
func getKeyFromString(keyString string) ([]byte, error) {
	return []byte(keyString), nil
}

// getKeyFromFile gets the encryption key from provided file.
func getKeyFromFile(keyFile string) ([]byte, error) {
	key, err := utils.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return key, nil
}

// getKeyFromPrompt gets the encryption key from user input.
func getKeyFromPrompt(encryptionRequired bool) ([]byte, error) {
	if encryptionRequired {
		bytePassword, err := readPassword("New password: ")
		if err != nil {
			return nil, fmt.Errorf("\nfailed to get key from prompt: %w", err)
		}

		checkPassword, err := readPassword("\nRepeat password: ")
		fmt.Println()
		if err != nil {
			return nil, fmt.Errorf("failed to get key from prompt: %w", err)
		}

		if !bytes.Equal(bytePassword, checkPassword) {
			return nil, errors.New("passwords don't match")
		}

		key := bytePassword

		return key, nil
	}
	key, err := readPassword("Password: ")
	fmt.Println()
	if err != nil {
		return nil, fmt.Errorf("failed to get key from prompt: %w", err)
	}

	return key, nil
}

// readPassword reads the user input from STDIN.
func readPassword(promt string) ([]byte, error) {
	fmt.Print(promt)

	stdin := int(syscall.Stdin)
	oldState, err := term.GetState(stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read password: %w", err)
	}
	defer term.Restore(stdin, oldState)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		for range sigch {
			term.Restore(stdin, oldState)
			fmt.Print("\nReceived interrupt signal.\n")
			os.Exit(1)
		}
	}()

	password, err := term.ReadPassword(stdin)
	if err != nil {
		return nil, err
	}

	return password, nil
}
