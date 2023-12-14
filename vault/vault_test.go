package vault

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	t.Run("test getKeyFromString", func(t *testing.T) {

		km := &KeyType{
			keyFile:            "",
			keyString:          "test_encrypt_string",
			encryptionRequired: false,
		}

		key, err := getKeyFromString(km.keyString)

		assert.Equal(t, len(key), len(key))
		assert.Nil(t, err)

	})

	t.Run("test getKeyFromFile", func(t *testing.T) {

		tempFile, _, cleanupTempFile := createTempFileWithContent(t, "password", "test_encrypt_string")
		defer cleanupTempFile()

		km := &KeyType{
			keyFile:            tempFile.Name(),
			keyString:          "",
			encryptionRequired: false,
		}

		key, err := getKeyFromFile(km.keyFile)

		assert.Equal(t, len(key), len(key))
		assert.Nil(t, err)
	})
}

func TestEncryptionDecryption(t *testing.T) {
	tempFile, tempContent, cleanupTempFile := createTempFileWithContent(t, "tempFile", "This is some text to encrypt!")
	passFile, passContent, cleanupPassFile := createTempFileWithContent(t, "passFile", "test_encrypt_string")

	defer cleanupTempFile()
	defer cleanupPassFile()

	km := &KeyType{
		keyFile:            passFile.Name(),
		keyString:          string(passContent),
		encryptionRequired: true,
	}

	t.Run("test encryption", func(t *testing.T) {
		err := Encrypt(tempFile.Name(), "", km.keyFile, "", true)
		assert.Nil(t, err, "Error should be nil.")
	})
	t.Run("test decryption", func(t *testing.T) {
		decryptedContent, err := Decrypt(tempFile.Name(), "", "", string(passContent), true)
		assert.Equal(t, tempContent, decryptedContent, "The contents should be equal.")
		assert.Nil(t, err, "Error should be nil.")
	})
}

func createTempFileWithContent(t *testing.T, file, text string) (*os.File, []byte, func()) {
	t.Helper()
	tempFile, err := os.CreateTemp("", file)
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}

	cleanup := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	content := []byte(text)
	writeContentToFile(t, tempFile, content)

	return tempFile, content, cleanup
}

func writeContentToFile(t *testing.T, file *os.File, content []byte) {
	t.Helper()
	_, err := file.Write(content)
	if err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}
}
