package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "account",
		Usage: "transfer",
		Action: func(c *cli.Context) error {
			fmt.Println("transfer")
			return nil
		},
	}
	app.Commands = []cli.Command{
		{
			Name:     "transfer",
			Aliases:  []string{"t"},
			Usage:    "transfer",
			Category: "transfer",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Run(os.Args)
}
