package cpchain_test

import (
	"context"
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
	wallet, err := clientOnTestnet.LoadWallet(keystorePath)
	if err != nil {
		t.Fatal(err)
	}
	expectAddr := "0xFD15C2932a60631222F7e6ffDdE7bDAB7237C2dC"
	if wallet.Addr().Hex() != expectAddr {
		t.Fatalf("expect %v to be %v", wallet.Addr().Hex(), expectAddr)
	}
	k, err := wallet.GetKey(password)
	if err != nil {
		t.Fatal(err)
	}
	_ = k
	// TODO validate private key
	// k.PrivateKey
}

func TestGetNonce(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, _ := clientOnTestnet.LoadWallet(keystorePath)

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

func TestSignTx(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := clientOnTestnet.LoadWallet(keystorePath)
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

	signedTx, err := wallet.SignTxWithPassword(password, tx, chainID)

	if err != nil {
		t.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)

	if err != nil {
		t.Fatal(err)
	}

}
