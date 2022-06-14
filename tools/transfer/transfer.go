package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	// "io/ioutil"
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
	"github.com/urfave/cli"
	"github.com/zgljl2012/slog"
)

// var (
// 	password string
//     passwordagian string
// )

const (
	Cpc = 1e18
)

func main() {
	app := &cli.App{
		Name:  "account",
		Usage: "transfer",
		Action: func(c *cli.Context) error {

			value := int64(1)
			endpoint := "https://civilian.testnet.cpchain.io"
			keystorePath := "e:/chengtcode/cpchain-golang-sdk/keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977"
			targetAddr := "0x4f5625efef254760301d2766c6cc98f05722963e"
			chainId := uint64(41)
			password := "test123456!"

			// value := c.Int64("value")
			// endpoint := c.String("endpoint")
			// keystorePath := c.String("keystore")
			// targetAddr := c.String("to")
			// chainId := c.Uint64("chainId")

			// if !c.IsSet("value") ||
			// 	!c.IsSet("endpoint") ||
			// 	!c.IsSet("to") ||
			// 	!c.IsSet("keystore") ||
			// 	!c.IsSet("chainId") {
			// 		fmt.Println("Need more parameter ! Check parameters with ./transfer -h please.")
			// 	return cli.NewExitError("Need more parameter ! Check parameters with ./transfer -h please. ", 1)
			// }
			// fpath, er := filepath.Abs("./keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977")
			fpath, er := filepath.Abs(keystorePath)
			if er != nil {
				fmt.Println(er)
			}
			// fmt.Println("please input your password") //TODO 改成其他形式的
			// fmt.Scanln(&password)
			// fmt.Println("please input your password again")
			// fmt.Scanln(&passwordagian)
			// if passwordagian == password {

			// } else {
			// 	return cli.NewExitError("ERROR: the password did not match the re-typed password", 1)
			// }
			password = "test123456!"

			client, err, privateKey, _, fromAddress, _, _ := cpchain.Connect(endpoint, fpath, password)

			nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
			if err != nil {
				slog.Fatal(err)
			}
			gasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				slog.Fatal(err)
			}

			to := common.HexToAddress(targetAddr)

			valueInCpc := new(big.Int).Mul(big.NewInt(value), big.NewInt(Cpc))

			msg := cpchain.CallMsg{From: fromAddress, To: &to, Value: valueInCpc, Data: nil}

			gasLimit, err := client.EstimateGas(context.Background(), msg)
			if err != nil {
				slog.Fatal(err)
			}

			tx := types.NewTransaction(nonce, to, valueInCpc, gasLimit, gasPrice, nil)

			chainID := big.NewInt(0).SetUint64(chainId)

			signedTx, err := types.SignTx(tx, types.NewCep1Signer(chainID), privateKey)

			if err != nil {
				slog.Fatal(err)
			}
			fmt.Println(tx)
			fmt.Println(signedTx)
			err = client.SendTransaction(context.Background(), signedTx)

			if err != nil {
				slog.Fatal(err)
			}
			return err
		},
	}

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:  "endpoint, ep",
	// 		Usage: "Endpoint to interact with",
	// 		Value: "http://localhost:8501",
	// 	},

	// 	cli.StringFlag{
	// 		Name:  "keystore, ks",
	// 		Usage: "Keystore dir path for from address,only 1 keystore file under the dir,path must explicit end with '/'",
	// 		Value: "./keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977",
	// 	},

	// 	cli.StringFlag{
	// 		Name:  "to",
	// 		Usage: "Recipient address",
	// 	},

	// 	cli.IntFlag{
	// 		Name:  "value",
	// 		Usage: "Value in cpc",
	// 	},

	// 	cli.IntFlag{
	// 		Name:  "chainId",
	// 		Usage: "chainId",
	// 	},
	// }

	app.Run(os.Args)
}
