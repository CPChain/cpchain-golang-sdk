package fusion_test

import (
	"math/big"
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
)

func TestWeb3(t *testing.T) {
	URL := "https://civilian.cpchain.io"
	provider, err := fusion.NewHttpProvider(URL)
	if err != nil {
		t.Fatal(err)
	}
	cf, err := fusion.NewWeb3(provider)
	if err != nil {
		t.Fatal(err)
	}
	rawBlock, err := cf.GetBlock("latest")
	if err != nil {
		t.Fatal(err)
	}
	block := rawBlock.(*fusion.FullBlock)
	if block.Number < 100 {
		t.Fatal("Block Failed")
	}
	rawBlock, err = cf.GetBlock(6435247)
	if err != nil {
		t.Fatal(err)
	}
	block = rawBlock.(*fusion.FullBlock)
	if block.Number != 6435247 {
		t.Fatal("Get block failed")
	}
	if block.Hash != "0xb19f8cde168407a7e2969f24cd35239f4beade33224852158010c288811b3e90" {
		t.Fatal("Block Hash is error")
	}
	txs := block.Transactions
	if len(txs) != 1 {
		t.Fatal("Transactions'cnt is error")
	}
	tx1 := block.Transactions[0]
	actual := big.NewInt(2.21 * fusion.Cpc)
	if actual.Cmp(tx1.Value) != 0 {
		t.Fatal("Actual value is error")
	}
}

func TestGetBalance(t *testing.T) {
	URL := "https://civilian.cpchain.io"
	provider, err := fusion.NewHttpProvider(URL)
	if err != nil {
		t.Fatal(err)
	}
	cf, err := fusion.NewWeb3(provider)
	if err != nil {
		t.Fatal(err)
	}
	r := cf.GetBalance("0xfe8c03415df612dc0e8c866283a4ed40277fa48b")
	cpc := r.Div(r, big.NewInt(fusion.Cpc))
	t.Logf("Balance at latest: %v", cpc)

	r, _ = cf.GetBalanceAt("0xfe8c03415df612dc0e8c866283a4ed40277fa48b", 6451549)
	cpc = r.Div(r, big.NewInt(fusion.Cpc))
	t.Logf("Balance at latest: %v", cpc)
}
