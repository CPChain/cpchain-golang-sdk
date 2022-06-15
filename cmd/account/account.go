package main

import (
	"fmt"
	"os"

	// "github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/urfave/cli"
)

var password string

var passwordfirst string
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
			Action: func(c *cli.Context) error {
				fmt.Println("please input your password")
				fmt.Scanln(&passwordfirst)
				fmt.Println("please input your password again")
				fmt.Scanln(&passwordagian)
				if passwordagian == passwordfirst {
					// err := cpchain.CreateKeystore(password)
					// err := nil
					// if err != nil {
					// 	fmt.Printf("%c[0;40;31m%s%c[0m", 0x1b, "ERROR: create new keystore failed", 0x1b)
					// }
				} else {
					fmt.Printf("%c[0;40;31m%s%c[0m", 0x1b, "ERROR: the password did not match the re-typed password", 0x1b)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
