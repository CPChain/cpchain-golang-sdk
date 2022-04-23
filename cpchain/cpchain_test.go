package cpchain_test

import (
	"cpchain-golang-sdk/cpchain"
	"math/big"
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	clientOnMainnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	blockNumberOnMainnet, err := clientOnMainnet.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if blockNumberOnMainnet == 0 {
		t.Fatal("BlockNumber is 0")
	}
	block1, err := clientOnMainnet.Block(1)
	if err != nil {
		t.Fatal(err)
	}
	if block1.Number != 1 {
		t.Fatal("BlockNumber is error")
	}
}

func TestGetBalance(t *testing.T) {
	clientOnMainnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	balanceOnMainnet := clientOnMainnet.BalanceOf("0x0a1ea332c4d3d457f17e0ada059f7275b3e2ea1e")
	if balanceOnMainnet.Cmp(big.NewInt(0)) == 0 {
		t.Fatal("Balance is 0")
	}
	t.Log(cpchain.WeiToCpc(balanceOnMainnet))
}
