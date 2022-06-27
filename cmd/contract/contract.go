package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	// "io/ioutil"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/tools"
	"github.com/urfave/cli"
	"github.com/zgljl2012/slog"
)

func main() {
	app := &cli.App{
		Name:  "contract",
		Usage: "deploy contract",
		Action: func(c *cli.Context) error {
			endpoint := c.String("endpoint")
			keystorePath := c.String("keystore")
			contractFilePath := c.String("contractfile")

			if !c.IsSet("endpoint") ||
				!c.IsSet("keystore") ||
				!c.IsSet("contractfile") {
				fmt.Println("Need more parameter ! Check parameters with ./transfer -h please.")
				return cli.NewExitError("Need more parameter ! Check parameters with ./transfer -h please. ", 1)
			}
			fpath, err := filepath.Abs(keystorePath)

			cfpath, err := filepath.Abs(contractFilePath)

			if err != nil {
				slog.Fatal(err)
			}
			password, err := tools.GetPassword("Please input your password:", false) //TODO 不显示星号
			network, err := cpchain.GetNetWork(endpoint)
			if err != nil {
				slog.Fatal(err)
			}
			clientOnTestnet, err := cpchain.NewCPChain(network) //TODO 根据endpoint
			if err != nil {
				slog.Fatal(err)
			}
			wallet, err := clientOnTestnet.LoadWallet(fpath)
			confirm := tools.AskForConfirmation("Are you sure to deploy this contract?")
			if !confirm {
				return nil
			}
			address, tx, err := wallet.DeployContractByFile(cfpath, password)
			time.Sleep(10 * time.Second)
			fmt.Println("Account:", address.Hex())
			fmt.Println("Tx hash:", tx.Hash().Hex())
			fmt.Println("Contract Depolyed!")
			return err
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "endpoint, ep",
			Usage: "Endpoint to interact with",
			Value: "http://localhost:8501",
		},

		cli.StringFlag{
			Name:  "keystore, ks",
			Usage: "Keystore file path for from address",
			Value: "./fixtures/keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977",
		},

		cli.StringFlag{
			Name:  "contractfile, cf",
			Usage: "contract file path for contract",
			Value: "./fixtures/contract/helloworld.js",
		},
	}

	app.Run(os.Args)
}
