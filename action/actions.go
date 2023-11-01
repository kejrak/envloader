package action

import (
	"fmt"

	"github.com/kejrak/envLoader/vault"
	"github.com/urfave/cli/v2"
)

func Encrypt(cCtx *cli.Context) error {
	err := vault.Encrypt(cCtx.String("file"),
		cCtx.String("output"),
		cCtx.String("key-file"),
		cCtx.String("key"))

	if cCtx.String("key") != "" {
		fmt.Println("Using --key via the CLI is insecure!")
	}

	if err != nil {
		fmt.Print("Encryption failed!\n")
		return fmt.Errorf("%v", err)
	}

	fmt.Print("Encryption successful.\n")

	return nil
}

func Decrypt(cCtx *cli.Context) error {
	_, err := vault.Decrypt(cCtx.String("file"),
		cCtx.String("output"),
		cCtx.String("key-file"),
		cCtx.String("key"))

	if cCtx.String("key") != "" {
		fmt.Println("Using --key via the CLI is insecure!")
	}

	if err != nil {
		fmt.Print("Decryption failed!\n")
		return fmt.Errorf("%v", err)
	}

	fmt.Print("Decryption successful.\n")

	return nil
}

func Load(cCtx *cli.Context) error {
	err := vault.Load(cCtx.String("file"),
		cCtx.String("binary"),
		cCtx.String("environment"),
		cCtx.String("key-file"),
		cCtx.String("key"))

	if cCtx.String("key") != "" {
		fmt.Println("Using --key via the CLI is insecure!")
	}

	if err != nil {
		fmt.Print("Loading environment variables failed!\n")
		return err
	}

	return nil
}
