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
	wallet, _ := clientOnTestnet.LoadWallet(keystorePath)

	tx, err := wallet.Transfer(password, targetAddr, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Tx hash: %v", tx.Hash().Hex())
}

func TestWalletDeploy(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, _ := clientOnTestnet.LoadWallet(keystorePath)

	address, tx, err := wallet.DeployContractByFile(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("address: %v", address.Hex())
	t.Logf("txhash: %v", tx.Hash().Hex())
	t.Logf("tx: %v", tx)
}
