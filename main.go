package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version = "1.0.0"
)

func main() {
	app := &cli.App{
		Name:        "envLoader",
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
