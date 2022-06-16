package cpchain_test

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
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

// 测试合约的事件
type CreateProductEvent struct {
	Id        cpchain.UInt256 `json:"ID"`
	Name      cpchain.String  `json:"name"`
	Extend    cpchain.String  `json:"extend"`
	Price     cpchain.UInt256 `json:"price"`
	Creator   cpchain.Address `json:"creator"`
	File_uri  cpchain.String  `json:"file_uri" rlp:"file_uri"`
	File_hash cpchain.String  `json:"file_hash"`
}

func TestEvents(t *testing.T) {
	client, err := cpchain.NewCPChain(cpchain.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	file, err := ioutil.ReadFile("../fixtures/product.json")
	if err != nil {
		t.Fatal(err)
	}
	address := "0x49F431A6bE97bd26bD416D6E6A0D3FAF3E3d5071"
	events, err := client.Contract(file, address).Events("CreateProduct",
		CreateProductEvent{},
		cpchain.WithEventsOptionsFromBlock(6712515))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Count:", len(events))
	for _, e := range events {
		args := e.Data.(*CreateProductEvent)
		t.Log(e.BlockNumber, args.Id, args.Name, args.Price, args.Extend, args.File_hash, args.File_uri, args.Creator.Hex())
		// check event name
		if e.Name != "CreateProduct" {
			t.Fatal("event name is error")
		}
	}
}

func TestCreateWallet(t *testing.T) {
	password := "123456"
	client, err := cpchain.NewCPChain(cpchain.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	path, err := ioutil.TempDir("e:/chengtcode/cpchain-golang-sdk/fixtures", "keystore")
	a, err := client.CreateWallet(path, password)
	if err != nil {
		t.Fatal(err)
	}
	w := client.LoadWallet(a.URL.Path)
	key, err := w.GetKey(password)
	if key.Address != a.Address {
		t.Fatal("account error")
	}
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll(path)
}
