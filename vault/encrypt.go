package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kejrak/envLoader/utils"
)

// Encrypt encrypts the specified file using AES-256-GCM encryption.
// It takes input file, output file (optional), key file, key string, and an inplace flag.
// If inplace is true, it overwrites the input file with the encrypted content.
// If output is specified, it writes the encrypted content to the output file.
func Encrypt(file, output, keyFile, keyString string, inplace bool) error {

	km := &KeyType{
		keyFile:            keyFile,
		keyString:          keyString,
		encryptionRequired: true,
	}

	checkEncrypted, err := utils.CheckEncryptedFile(file)
	if err != nil {
		fmt.Printf("can't check file!")
		return err
	}

	if checkEncrypted {
		return fmt.Errorf("file is already encrypted")
	}

	key, err := km.getKey()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	key, salt, err := km.deriveKey(key, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	cipherText, err := encryption(file, key, salt)
	if err != nil {
		return err
	}

	if output != "" {
		utils.WriteFile(output, cipherText)
	} else if inplace {
		utils.WriteFile(file, cipherText)
	} else {
		fmt.Fprint(os.Stdout, string(cipherText))
	}

	return nil
}

// encryption performs AES-GCM encryption on the content of the specified file.
// It returns the encrypted content along with a special header.
func encryption(file string, key, salt []byte) ([]byte, error) {

	plainText, err := os.ReadFile(file)

	if err != nil {
		log.Printf("read file err: %v", err)
		return nil, err
	}

	secretKey := []byte(key)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		log.Printf("cipher err: %v", err)
		return nil, err
	}

	cipherText, err := gcmEncrypt(plainText, block)
	if err != nil {
		log.Printf("cannot encrypt: %v", err)
		return nil, err
	}

	cipherText = append(cipherText, salt...)
	encodedCiphertext := base64.StdEncoding.EncodeToString(cipherText)

	result := []byte("!envLoader | AES-256\n" + encodedCiphertext)
	if err != nil {
		log.Printf("write file err: %v", err.Error())
	}

	return []byte(utils.WrapBytes(result)), nil
}

// gcmEncrypt performs AES-GCM encryption on the given plaintext using the provided block.
func gcmEncrypt(plainText []byte, block cipher.Block) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("cipher GCM err: %v", err)
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.Printf("nonce  err: %v", err)
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plainText, nil), nil
}
