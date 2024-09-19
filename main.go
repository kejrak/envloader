package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kejrak/envloader/utils"
	"github.com/urfave/cli/v2"
)

var (
	version = ""
)

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s	
SHELL TYPE: 
	 %s	Type of shell, that runs script (default '/bin/sh')
`, cli.AppHelpTemplate, utils.GetEnv("ENVLOADER_SHELL_TYPE", "/bin/sh"))

	app := &cli.App{
		Name:        "envloader",
		Usage:       "environment cli tool",
		Version:     version,
		Description: "A cli tool to inject variables from encrypted / decrypted file into binary.",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Jan Kejř",
				Email: "jan.kejr@centrum.cz",
			},
		},
		UseShortOptionHandling: true,
		Commands:               commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
