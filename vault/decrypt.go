package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kejrak/envLoader/utils"
)

// Decrypt takes an encrypted file, decrypts its content using AES-256 algorithm,
// and writes the decrypted content to an output file or stdout.
// If inplace is true, the original file is overwritten with the decrypted content.
func Decrypt(file, output, keyFile, keyString string, inplace bool) ([]byte, error) {

	km := &KeyType{
		keyFile:            keyFile,
		keyString:          keyString,
		encryptionRequired: false,
	}

	checkEncrypted, err := utils.CheckEncryptedFile(file)
	if err != nil {
		return nil, err
	}

	if !checkEncrypted {
		return nil, errors.New("file is already decrypted")
	}

	key, err := km.getKey()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	plainText, err := decryption(file, key, km)

	if err != nil {
		return nil, err
	}

	if output != "" {
		utils.WriteFile(output, plainText)
	} else if inplace {
		utils.WriteFile(file, plainText)
	} else {
		fmt.Fprint(os.Stdout, string(plainText))
	}

	return plainText, nil
}

// decryption performs AES-GCM decryption on the content of the specified file.
// It returns the decrypted content.
func decryption(file string, key []byte, keyManager *KeyType) ([]byte, error) {

	cipherText, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	cipherText = []byte(strings.TrimPrefix(string(cipherText), "!envLoader | AES-256"))

	cipherText, err = base64.StdEncoding.DecodeString(strings.TrimSpace(string(cipherText)))
	if err != nil {
		return nil, err
	}

	salt, cipherText := cipherText[len(cipherText)-32:], cipherText[:len(cipherText)-32]
	key, _, err = keyManager.deriveKey(key, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	plainText, err := gcmDecrypt(cipherText, block)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return plainText, nil
}

// gcmDecrypt performs AES-GCM decryption on the given ciphertext using the provided block.
// It returns the decrypted plaintext.
func gcmDecrypt(cipherText []byte, block cipher.Block) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
