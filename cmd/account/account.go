package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/tools"
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
				password, err := tools.GetPassword("Please input your password:", true)
				if err != nil {
					return nil
				}
				client, err := cpchain.NewCPChain(cpchain.Testnet)
				if err != nil {
					slog.Fatal(err)
				}
				a, err := client.CreateAccount(keystorePath, password)
				if err != nil {
					slog.Fatal(err)
				}
				fmt.Println("Account:", a.Address.Hex())
				fmt.Println("Path:", a.URL.Path)
				fmt.Println("Create account successful")
				return nil
			},
		},
	}
	app.Run(os.Args)
}
