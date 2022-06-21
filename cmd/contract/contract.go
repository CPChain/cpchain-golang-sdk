package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	// "io/ioutil"
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/cmd/contract/utils"
	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/contract"
	"github.com/urfave/cli"
	"github.com/zgljl2012/slog"
)

var (
	password      string
	passwordagian string
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

			Abi, Bin, err := utils.ReadContract(cfpath)
			if err != nil {
				slog.Fatal(err)
			}
			fmt.Println("please input your password") //TODO 改成其他形式的
			fmt.Scanln(&password)
			fmt.Println("please input your password again")
			fmt.Scanln(&passwordagian)
			if passwordagian == password {

			} else {
				return cli.NewExitError("ERROR: the password did not match the re-typed password", 1)
			}
			network, err := utils.GetNetWork(endpoint)
			if err != nil {
				slog.Fatal(err)
			}
			clientOnTestnet, err := cpchain.NewCPChain(network) //TODO 根据endpoint
			if err != nil {
				slog.Fatal(err)
			}
			wallet := clientOnTestnet.LoadWallet(fpath)

			backend, err := clientOnTestnet.Backend()
			if err != nil {
				slog.Fatal(err)
			}

			Key, err := wallet.GetKey(password)
			if err != nil {
				slog.Fatal(err)
			}
			fromAddr := wallet.Addr()

			nonce, err := backend.PendingNonceAt(context.Background(), fromAddr)
			if err != nil {
				slog.Fatal(err)
			}

			auth := contract.NewTransactor(Key.PrivateKey, new(big.Int).SetUint64(nonce))

			fmt.Println("start depoly contract")
			_, _, _, err = clientOnTestnet.DeployContract(Abi, Bin, auth)
			if err != nil {
				slog.Fatal(err)
			}
			fmt.Println("success")
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
