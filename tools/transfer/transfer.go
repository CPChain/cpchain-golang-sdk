package main

import (
	"fmt"
	"os"
	"path/filepath"

	// "io/ioutil"
	"github.com/CPChain/cpchain-golang-sdk/cpchain"
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
				fpath, er := filepath.Abs("./keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977")
				if er != nil {
					fmt.Println(er)
				}
				// content, err := ioutil.ReadFile(fpath)
				// fmt.Println(string(content))
				privateKey, publicKey, _, _, _, err := cpchain.DecryptKeystore(fpath, "test123456!")
				fmt.Println(privateKey)
				fmt.Println(publicKey)

				// _, err := os.Open(fpath) // For read access.
				if err != nil {
					fmt.Println(err)
				}
				return err
			},
		},
	}

	app.Run(os.Args)
}
