package main

import (
	"fmt"
	"os"
	"path/filepath"

	// "io/ioutil"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/tools"
	"github.com/urfave/cli"
	"github.com/zgljl2012/slog"
)

const (
	Cpc = 1e18
)

func main() {
	app := &cli.App{
		Name:  "transfer",
		Usage: "transfer",
		Action: func(c *cli.Context) error {
			value := c.Int64("value")
			endpoint := c.String("endpoint")
			keystorePath := c.String("keystore")
			targetAddr := c.String("to")

			if !c.IsSet("value") ||
				!c.IsSet("endpoint") ||
				!c.IsSet("to") ||
				!c.IsSet("keystore") {
				fmt.Println("Need more parameter ! Check parameters with ./transfer -h please.")
				return cli.NewExitError("Need more parameter ! Check parameters with ./transfer -h please. ", 1)
			}
			fpath, er := filepath.Abs(keystorePath)
			if er != nil {
				fmt.Println(er)
			}

			password, err := tools.GetPassword("Please input your password:", true)
			if err != nil {
				return cli.NewExitError("ERROR:", 1)
			}
			network, err := cpchain.GetNetWork(endpoint)
			if err != nil {
				slog.Fatal(err)
			}
			clientOnTestnet, err := cpchain.NewCPChain(network) //TODO 根据endpoint
			if err != nil {
				slog.Fatal(err)
			}
			wallet := clientOnTestnet.LoadWallet(fpath)

			err = wallet.Transfer(targetAddr, password, value)
			fmt.Println("success")
			return err
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "endpoint, ep",
			Usage: "Endpoint to interact with",
			Value: "https://civilian.testnet.cpchain.io",
		},

		cli.StringFlag{
			Name:  "keystore, ks",
			Usage: "Keystore dir path for from address,only 1 keystore file under the dir,path must explicit end with '/'",
			Value: "./keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977",
		},

		cli.StringFlag{
			Name:  "to",
			Usage: "Recipient address",
		},

		cli.IntFlag{
			Name:  "value",
			Usage: "Value in cpc",
		},
	}

	app.Run(os.Args)
}
