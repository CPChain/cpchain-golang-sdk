package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/urfave/cli"
	"github.com/zgljl2012/slog"
)

var password string

var passwordagian string

func main() {
	app := &cli.App{
		Name:  "account",
		Usage: "create account",
		Action: func(c *cli.Context) error {
			fmt.Println("create")
			return nil
		},
	}
	app.Commands = []cli.Command{
		{
			Name:     "new",
			Aliases:  []string{"n"},
			Usage:    "create a new account",
			Category: "account",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keystorepath , kp",
					Usage: "The path where the generated files are stored",
					Value: "./keystore",
				}},
			Action: func(c *cli.Context) error {
				keystorePath := c.String("keystorepath")
				fmt.Println(keystorePath)
				if !c.IsSet("keystorepath") {
					fmt.Println("Need more parameter ! Check parameters with ./account -h please.")
					return cli.NewExitError("Need more parameter ! Check parameters with ./account -h please. ", 1)
				}
				keystorePath, err := filepath.Abs(keystorePath)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("please input your password")
				fmt.Scanln(&password)
				fmt.Println("please input your password again")
				fmt.Scanln(&passwordagian)
				if passwordagian == password {
					client, err := cpchain.NewCPChain(cpchain.Testnet)
					if err != nil {
						slog.Fatal(err)
					}
					a, err := client.CreateAccount("e:/chengtcode/cpchain-golang-sdk/fixtures/keystore", password)
					if err != nil {
						slog.Fatal(err)
					}
					fmt.Println("account:", a.Address.Hex())
					fmt.Println("path:", a.URL.Path)
				} else {
					fmt.Printf("%c[0;40;31m%s%c[0m", 0x1b, "ERROR: the password did not match the re-typed password", 0x1b)
				}
				fmt.Println("success")
				return nil
			},
		},
	}
	app.Run(os.Args)
}
