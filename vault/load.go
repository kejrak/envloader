package vault

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kejrak/envLoader/utils"
	"gopkg.in/ini.v1"
)

// Load reads environment configuration from a file, decrypts if necessary,
// and loads it into a binary by executing a command based on environmental value.
func Load(file, binary, envType, keyFile, keyString string) error {
	plainText, err := readPlainText(file, keyFile, keyString)
	if err != nil {
		return err
	}

	cfg, err := ini.Load(plainText)
	if err != nil {
		fmt.Printf("fail to load file: %v\n", err)
	}

	globalSection := cfg.Section("")
	section := cfg.Section(envType)

	if len(section.Keys()) != 0 {

		err := loadToBinary(envType, binary, globalSection, section)
		if err != nil {
			fmt.Printf("fail to load env variables to binary: %v\n", err)
		}
		return nil

	} else {
		fmt.Print("environment configuration doesn't exist!\n")
		return nil
	}

}

// readPlainText reads the plaintext content of the file, decrypting if necessary.
func readPlainText(file, keyFile, keyString string) ([]byte, error) {
	var plainText []byte

	km := &KeyType{
		keyFile:            keyFile,
		keyString:          keyString,
		encryptionRequired: false,
	}

	checkEncrypted, err := utils.CheckEncryptedFile(file)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	if !checkEncrypted {
		fmt.Print("file is decrypted!\n")

		plainText, err = utils.ReadFile(file)
		if err != nil {
			fmt.Printf("%v", err)
			return nil, err
		}

		return plainText, nil

	} else {
		key, err := km.getKey()
		if err != nil {
			fmt.Printf("%v", err)
			return nil, err
		}

		plainText, err = decryption(file, key, km)
		if err != nil {
			fmt.Printf("%v", err)
			return nil, err
		}

		return plainText, nil
	}
}

// loadToBinary loads environment variables into a binary by executing a command based on environmental value.
func loadToBinary(envType, binary string, globalSection, section *ini.Section) error {
	if _, err := os.Stat(binary); os.IsNotExist(err) {
		fmt.Printf("binary file doesn't exist!\n")
		return nil
	}

	fmt.Printf("configuration: %s\n", envType)

	cmd := exec.Command("/bin/sh", binary)

	appendKeysToCmdEnv(globalSection.Keys(), cmd)
	appendKeysToCmdEnv(section.Keys(), cmd)

	return cmd.Run()
}

// appendKeysToCmdEnv appends INI keys to the environment of a command based on environmental value.
func appendKeysToCmdEnv(keys []*ini.Key, cmd *exec.Cmd) {
	for _, key := range keys {
		value := fmt.Sprintf("%s=%s", strings.ToUpper(key.Name()), key.Value())
		cmd.Env = append(cmd.Env, value)
	}
}
