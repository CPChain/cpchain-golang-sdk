package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
	"github.com/CPChain/cpchain-golang-sdk/tools"
	"github.com/urfave/cli"
	"github.com/zgljl2012/slog"
)

var password string

var passwordagian string

func main() {
	app := &cli.App{
		Name:  "contract",
		Usage: "depoly contract and call contract function",
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	app.Commands = []cli.Command{
		{
			Name:     "deploy",
			Aliases:  []string{"d"},
			Usage:    "deploy contract",
			Category: "contract",
			Flags: []cli.Flag{
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
					Value: "./fixtures/contract/Hello.js",
				},
			},
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
				wallet, err := clientOnTestnet.LoadWallet(fpath, password)
				confirm := tools.AskForConfirmation("Are you sure to deploy this contract?")
				if !confirm {
					return nil
				}
				address, tx, err := clientOnTestnet.DeployContractByFile(cfpath, wallet)
				time.Sleep(10 * time.Second)
				fmt.Println("Account:", address.Hex())
				fmt.Println("Tx hash:", tx.Hash().Hex())
				receipt, err := clientOnTestnet.ReceiptByTx(tx)
				if err != nil {
					slog.Fatal(err)
				}
				if receipt.Status == types.ReceiptStatusSuccessful {
					fmt.Println("confirm transaction success")
				} else {
					fmt.Println("confirm transaction failed", "status", receipt.Status,
						"receipt.TxHash", receipt.TxHash)
				}
				fmt.Println("Contract Depolyed!")
				return err
			},
		},
		{
			Name:     "call",
			Aliases:  []string{"c"},
			Usage:    "call contract function",
			Category: "contract",
			Flags: []cli.Flag{
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
					Value: "./fixtures/contract/Hello.json",
				},

				cli.StringFlag{
					Name:  "contractaddr, ca",
					Usage: "address for contract",
					Value: "0x0",
				},

				cli.StringFlag{
					Name:  "function, f",
					Usage: "function of contract",
					Value: "",
				},
				cli.StringSliceFlag{
					Name:  "params, p",
					Usage: "params of contract function",
					Value: nil,
				},
			},
			Action: func(c *cli.Context) error {
				endpoint := c.String("endpoint")
				keystorePath := c.String("keystore")
				contractFilePath := c.String("contractfile")
				contractAddr := c.String("contractaddr")
				method := c.String("function")
				params := c.StringSlice("params")

				if params != nil {
					fmt.Println("get params")
				}

				// convertParmas(params)
				if !c.IsSet("endpoint") ||
					!c.IsSet("keystore") ||
					!c.IsSet("contractfile") ||
					!c.IsSet("contractaddr") ||
					!c.IsSet("function") {
					fmt.Println("Need more parameter ! Check parameters with ./contract -h please.")
					return cli.NewExitError("Need more parameter ! Check parameters with ./contract -h please. ", 1)
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
				wallet, err := clientOnTestnet.LoadWallet(fpath, password)
				confirm := tools.AskForConfirmation("Are you sure to call this function?")
				if !confirm {
					return nil
				}
				abi, _, err := cpchain.ReadContract(cfpath)
				fmt.Println(abi, contractAddr, wallet.Addr(), method, network.ChainId)
				tx, err := clientOnTestnet.Contract([]byte(abi), contractAddr).Call(wallet, network.ChainId, method, params...)
				fmt.Println("Tx hash:", tx.Hash().Hex())
				fmt.Println("Please wait 10 seconds")
				receipt, err := clientOnTestnet.ReceiptByTx(tx)
				if err != nil {
					slog.Fatal(err)
				}
				if receipt.Status == types.ReceiptStatusSuccessful {
					fmt.Println("confirm transaction success")
				} else {
					fmt.Println("confirm transaction failed", "status", receipt.Status,
						"receipt.TxHash", receipt.TxHash)
				}
				return err
			},
		},
		{
			Name:     "view",
			Aliases:  []string{"c"},
			Usage:    "call contract function",
			Category: "contract",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "endpoint, ep",
					Usage: "Endpoint to interact with",
					Value: "http://localhost:8501",
				},

				cli.StringFlag{
					Name:  "contractfile, cf",
					Usage: "contract file path for contract",
					Value: "./fixtures/contract/Hello.json",
				},

				cli.StringFlag{
					Name:  "contractaddr, ca",
					Usage: "address for contract",
					Value: "0x0",
				},

				cli.StringFlag{
					Name:  "method, m",
					Usage: "method of contract",
					Value: "",
				},
				cli.StringSliceFlag{
					Name:  "params, p",
					Usage: "params of contract function",
					Value: nil,
				},
			},
			Action: func(c *cli.Context) error {
				endpoint := c.String("endpoint")
				contractFilePath := c.String("contractfile")
				contractAddr := c.String("contractaddr")
				method := c.String("method")
				params := c.StringSlice("params")

				// convertParmas(params)
				if !c.IsSet("endpoint") ||
					!c.IsSet("contractfile") ||
					!c.IsSet("contractaddr") ||
					!c.IsSet("method") {
					fmt.Println("Need more parameter ! Check parameters with ./contract -h please.")
					return cli.NewExitError("Need more parameter ! Check parameters with ./contract -h please. ", 1)
				}
				cfpath, err := filepath.Abs(contractFilePath)
				if err != nil {
					slog.Fatal(err)
				}
				network, err := cpchain.GetNetWork(endpoint)
				if err != nil {
					slog.Fatal(err)
				}
				clientOnTestnet, err := cpchain.NewCPChain(network) //TODO 根据endpoint
				if err != nil {
					slog.Fatal(err)
				}
				abi, _, err := cpchain.ReadContract(cfpath)
				result, err := clientOnTestnet.Contract([]byte(abi), contractAddr).View(method, params...)
				fmt.Print(result)
				return err
			},
		},
	}
	app.Run(os.Args)
}

// func convertParmas(p []string) []interface{} {
// 	var convertedParams []interface{}
// 	for k, v := range p {
// 		switch {
// 		case len(v) == 42 && v[0:2] == "0x": // 地址类型
// 			convertedParams[k] = common.HexToAddress(v)
// 		case v == "true":
// 			convertedParams[k] = true // 布尔类型
// 		case v == "false":
// 			convertedParams[k] = false
// 		case IsInt(v):
// 			// convertedParams[k] = strconv.ParseInt(v, 10, 64) //TODO 转换成uint256
// 		default:
// 			convertedParams[k] = v
// 		}
// 	}
// 	return convertedParams
// }

// func IsInt(p string) bool {
// 	if len(p) >= 2 && string(p[0]) == "0" {
// 		return false
// 	} else {
// 		patterns := regexp.MustCompile(`\d{1,}`)
// 		result := patterns.FindString(p)
// 		if result == p {
// 			return true
// 		}
// 		return false
// 	}
// }
