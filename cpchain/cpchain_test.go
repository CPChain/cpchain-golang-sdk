package cpchain_test

import (
	"cpchain-golang-sdk/cpchain"
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	clientOnMainnet, err := cpchain.NewCPChain(cpchain.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	blockNumberOnMainnet, err := clientOnMainnet.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blockNumberOnMainnet)
}
