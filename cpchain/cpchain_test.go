package cpchain_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/internal/cpcclient"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
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

func TestCreateAccount(t *testing.T) {
	password := "123456"
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	path, err := ioutil.TempDir(os.TempDir(), "keystore")
	if err != nil {
		t.Fatal(err)
	}
	a, err := client.CreateAccount(path, password)
	if err != nil {
		t.Fatal(err)
	}
	w, err := client.LoadWallet(a.URL.Path, password)
	if err != nil {
		t.Fatal(err)
	}
	key := w.Key()
	if key.Address != a.Address {
		t.Fatal("account error")
	}
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll(path)
}

const Abi = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

const Bin = `0x6080604052348015600f57600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550603f80605d6000396000f3fe6080604052600080fdfea2646970667358221220cc46356d887799b33b3ca82fcf610da45d06ecf8fa0e763740abfbd51f6898ff64736f6c634300080a0033`

func TestReadContract(t *testing.T) {
	abi, bin, err := cpchain.ReadContract(cfpath)
	t.Log(abi)
	t.Log(Abi)
	if err != nil {
		t.Fatal(err)
	}
	if abi != Abi {
		t.Fatal("abi! = Abi")
	}
	if bin != Bin {
		t.Fatal("bin != Bin")
	}
}

func TestContractDeploy(t *testing.T) {
	abi, bin, err := cpchain.ReadContract("../fixtures/contract/Hello.json")
	if err != nil {
		t.Fatal(err)
	}
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := client.LoadWallet(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}
	address, tx, err := client.DeployContract(abi, bin, wallet)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("address:%v", address.Hex())
	t.Logf("Tx hash: %v", tx.Hash().Hex())
}

// const helloaddress = "0x63c7AdaA9EEf9f89cbB8960AFAdca50782d259f4"
// const helloaddress = "0x9865Cb52E790e30E115a435d23683e2078a3ABB7"
const helloaddress = "0xfD44A7aEFaDfa872Ade30EBE152Fc37E6977fe70"

// 0x63c7AdaA9EEf9f89cbB8960AFAdca50782d259f4
var number []byte

var (
	Address, _ = abi.NewType("address")
)

func TestContractTransact(t *testing.T) {
	Abi, _, err := cpchain.ReadContract("../fixtures/contract/Hello.json")
	if err != nil {
		t.Fatal(err)
	}
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := client.LoadWallet(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}
	contracthello := client.Contract([]byte(Abi), helloaddress)
	tx, err := contracthello.Call(wallet, 41, "helloToEveryOne")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tx hash: %v", tx.Hash().Hex())
	receipt, err := client.ReceiptByTx(tx)

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

func TestU256(t *testing.T) {
	number = abi.U256(big.NewInt(4))
	fmt.Printf("value: %d\n", number)
	fmt.Printf("number of bytes: %d", len(number))
}

type HelloToSomeOne struct {
	Target cpchain.Address `json:"target"`
}

func TestContractEvent(t *testing.T) {
	Abi, _, err := cpchain.ReadContract("../fixtures/contract/Hello.json")
	if err != nil {
		t.Fatal(err)
	}
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	contracthello := client.Contract([]byte(Abi), helloaddress)
	events, err := contracthello.Events("HelloToSomeOne", HelloToSomeOne{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Count:", len(events))
	for _, e := range events {
		args := e.Data.(*HelloToSomeOne)
		t.Log(e.BlockNumber, args.Target.Hex())
		// check event name
		if e.Name != "HelloToSomeOne" {
			t.Fatal("event name is error")
		}
	}
}

func TestContractView(t *testing.T) {
	Abi, _, err := cpchain.ReadContract("../fixtures/contract/Hello.json")
	if err != nil {
		t.Fatal(err)
	}
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	contracthello := client.Contract([]byte(Abi), helloaddress)

	var hellotime = big.NewInt(0)
	// var hellotime *types.Receipt
	err = contracthello.View(&hellotime, "hellotime")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hellotime)
	if hellotime == nil || hellotime.Uint64() == 0 {
		t.Error("error")
	}
}

func TestReceiptByTx(t *testing.T) {
	client, err := cpcclient.Dial(cpchain.Testnet.JsonRpcUrl)
	rep, err := client.TransactionReceipt(context.Background(), common.HexToHash("0xa0aa285751bc0a1a9885705f1f24fd3ce2398b7a2377b952b78138d97741eaea"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rep)
}
