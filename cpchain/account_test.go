package cpchain_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/internal/cpcclient"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

var (
	value        = int64(1)
	endpoint     = "https://civilian.testnet.cpchain.io"
	keystorePath = "../fixtures/keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977"
	targetAddr   = "0x4f5625efef254760301d2766c6cc98f05722963e"
	chainId      = uint64(41)
	password     = "test123456!"
	cfpath       = "../fixtures/contract/helloworld.json"
)

const (
	Cpc = 1e18
)

func TestGetKey(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := clientOnTestnet.LoadWallet(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}
	expectAddr := "0xFD15C2932a60631222F7e6ffDdE7bDAB7237C2dC"
	if wallet.Addr().Hex() != expectAddr {
		t.Fatalf("expect %v to be %v", wallet.Addr().Hex(), expectAddr)
	}
	k := wallet.Key()
	_ = k
	// TODO validate private key
	// k.PrivateKey
}

func TestGetNonce(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, _ := clientOnTestnet.LoadWallet(keystorePath, password)

	client, err := cpcclient.Dial(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	fromAddr := wallet.Addr()

	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce:", nonce)
}

// ---- {"jsonrpc":"2.0","id":18,"result":{"blockHash":"0x7911ada70ca3187c5c63383e42f751f4941875d8715c464872d59cae33e244b7","blockNumber":"0x250aaa","contractAddress":null,"cumulativeGasUsed":"0x6d17","from":"0xfd15c2932a60631222f7e6ffdde7bdab7237c2dc","gasUsed":"0x6d17","logs":[{"address":"0xfd44a7aefadfa872ade30ebe152fc37e6977fe70","topics":["0x845d757b1759e3b909865aad71e093a4c4649a414515ee96c3ebb2d0bd18ed73"],"data":"0x0000000000000000000000000000000000000000000000000000000000000012","blockNumber":"0x250aaa","transactionHash":"0xbcd09ead48f7db0f9485cc1b6a7b6afb6d8a8dc72ea02d903d7aec90048344c2","transactionIndex":"0x0","blockHash":"0x7911ada70ca3187c5c63383e42f751f4941875d8715c464872d59cae33e244b7","logIndex":"0x0","removed":false}],"logsBloom":"0x000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000
// 00000004000000000000000000000020000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":"0xfd44a7aefadfa872ade30ebe152fc37e6977fe70","transactionHash":"0xbcd09ead48f7db0f9485cc1b6a7b6afb6d8a8dc72ea02d903d7aec90048344c2","transactionIndex":"0x0"}}
// ----+++++ {"blockHash":"0x7911ada70ca3187c5c63383e42f751f4941875d8715c464872d59cae33e244b7","blockNumber":"0x250aaa","contractAddress":null,"cumulativeGasUsed":"0x6d17","from":"0xfd15c2932a60631222f7e6ffdde7bdab7237c2dc","gasUsed":"0x6d17","logs":[{"address":"0xfd44a7aefadfa872ade30ebe152fc37e6977fe70","topics":["0x845d757b1759e3b909865aad71e093a4c4649a414515ee96c3ebb2d0bd18ed73"],"data":"0x0000000000000000000000000000000000000000000000000000000000000012","blockNumber":"0x250aaa","transactionHash":"0xbcd09ead48f7db0f9485cc1b6a7b6afb6d8a8dc72ea02d903d7aec90048344c2","transactionIndex":"0x0","blockHash":"0x7911ada70ca3187c5c63383e42f751f4941875d8715c464872d59cae33e244b7","logIndex":"0x0","removed":false}],"logsBloom":"0x00000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000004000000000000000000000
// 020000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":"0xfd44a7aefadfa872ade30ebe152fc37e6977fe70","transactionHash":"0xbcd09ead48f7db0f9485cc1b6a7b6afb6d8a8dc72ea02d903d7aec90048344c2","transactionIndex":"0x0"}

func TestSignTx(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := clientOnTestnet.LoadWallet(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}

	client, err := cpcclient.Dial(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	fromAddr := wallet.Addr()

	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		t.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	to := common.HexToAddress(targetAddr)

	valueInCpc := new(big.Int).Mul(big.NewInt(value), big.NewInt(Cpc))

	msg := cpcclient.CallMsg{From: fromAddr, To: &to, Value: valueInCpc, Data: nil}

	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}

	tx := types.NewTransaction(nonce, to, valueInCpc, gasLimit, gasPrice, nil)

	chainID := big.NewInt(0).SetUint64(chainId)

	signedTx, err := wallet.SignTx(tx, chainID)

	if err != nil {
		t.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(signedTx.Hash().Hex())
	// receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	receipt, err := clientOnTestnet.ReceiptByTx(signedTx)

	if err != nil {
		t.Fatalf("failed to waitMined tx:%v", err)
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		t.Log("confirm transaction success")
	} else {
		t.Error("confirm transaction failed", "status", receipt.Status,
			"receipt.TxHash", receipt.TxHash)
	}
}

// test receipt by txhash
func TestReceipt(t *testing.T) {
	client, err := cpcclient.Dial(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	hash_hex := "0x2c65bde3b32cd3bb05740330f3b5cd25455e2f81edae40c524b59f670b315998"
	hash := common.HexToHash(hash_hex)
	t.Log(hash)
	t.Log(hash.Hex())
	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(receipt.TxHash)

}
