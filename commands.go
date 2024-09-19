package main

import (
	"github.com/kejrak/envloader/action"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	{
		Name:  "encrypt",
		Usage: "encrypt file with key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "file to encrypt",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "key file with password",
			},
			&cli.StringFlag{
				Name:  "key-file",
				Usage: "key file with password",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "write output to file",
			},
			&cli.BoolFlag{
				Name:    "in-place",
				Aliases: []string{"i"},
				Usage:   "write output back to the same file instead of stdout",
			},
		},
		Action: action.Encrypt,
	},
	{
		Name:  "decrypt",
		Usage: "decrypt file with key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "file to decrypt",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "key file with password",
			},
			&cli.StringFlag{
				Name:  "key-file",
				Usage: "key file with password",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output file",
			},
			&cli.BoolFlag{
				Name:    "in-place",
				Aliases: []string{"i"},
				Usage:   "write output back to the same file instead of stdout",
			},
		},
		Action: action.Decrypt,
	},
	{
		Name:  "load",
		Usage: "load variables into binary file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "file to load",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "key file with password",
			},
			&cli.StringFlag{
				Name:  "key-file",
				Usage: "key file with password",
			},
			&cli.StringFlag{
				Name:    "binary",
				Aliases: []string{"b"},
				Usage:   "binary file",
			},
			&cli.StringFlag{
				Name:     "environment",
				Aliases:  []string{"e"},
				Usage:    "type of environment",
				Required: true,
			},
		},
		Action: action.Load,
	},
}
