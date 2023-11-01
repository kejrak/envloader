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

func Encrypt(file, output, keyFile, keyString string) error {

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

	cipherText, err := encryption(file, key)
	if err != nil {
		return err
	}

	if output != "" {
		utils.WriteFile(output, cipherText)
	} else {
		utils.WriteFile(file, cipherText)
	}

	return nil

}

func encryption(file, key string) ([]byte, error) {

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

	encodedCiphertext := base64.StdEncoding.EncodeToString(cipherText)

	result := []byte("!envLoader | AES-256\n" + encodedCiphertext)
	if err != nil {
		log.Printf("write file err: %v", err.Error())
	}

	cipherText = []byte(utils.WrapBytes(result))

	return cipherText, nil

}

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

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	return cipherText, nil

}
