package cpchain_test

import (
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
)

func TestWalletTransfer(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, _ := clientOnTestnet.LoadWallet(keystorePath, password)

	tx, err := wallet.Transfer(targetAddr, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Tx hash: %v", tx.Hash().Hex())
}
