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
	wallet := clientOnTestnet.LoadWallet(keystorePath)

	err = wallet.Transfer(password, targetAddr, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWalletDeploy(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet := clientOnTestnet.LoadWallet(keystorePath)

	err = wallet.DeployContractByFile(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}
}
