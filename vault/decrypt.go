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

func Decrypt(file, output, keyFile, keyString string) ([]byte, error) {

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
		return nil, err
	}

	plainText, err := decryption(file, key)

	if err != nil {
		return nil, err
	}

	if output != "" {
		utils.WriteFile(output, plainText)
	} else {
		utils.WriteFile(file, plainText)
	}

	return plainText, nil
}

func decryption(file, key string) ([]byte, error) {

	cipherText, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	cipherText = []byte(strings.TrimPrefix(string(cipherText), "!envLoader | AES-256"))

	cipherText, err = base64.StdEncoding.DecodeString(strings.TrimSpace(string(cipherText)))
	if err != nil {
		return nil, err
	}

	secretKey := []byte(key)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	result, err := gcmDecrypt(cipherText, block)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return result, nil

}

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
