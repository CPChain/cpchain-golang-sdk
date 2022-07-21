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

// 测试创建新钱包地址
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

// 测试从文件中读取contract abi 和bin
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

// 测试合约部署
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

// 测试call 合约的 function，会发送交易
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

	// ABI,err := abi.JSON(strings.NewReader(Abi))

	tx, err := contracthello.Call(wallet, 41, "helloToSomeOne", int64(0), targetAddr)
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

// 测试call 合约的 function，会发送交易(自动转换参数类型)
func TestContractTransactConvert(t *testing.T) {
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

	tx, err := contracthello.Call(wallet, 41, "helloToSomeOne", int64(0), targetAddr)
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
	// var hellotime *types.Receipt
	result, err := contracthello.View("hellotime")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
	if result == nil {
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

func Test1(t *testing.T) { //部署合约
	abi, bin, err := cpchain.ReadContract("../fixtures/contract/AirDrop.json")
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

func Test2(t *testing.T) { //给合约充值
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := client.LoadWallet(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}

	Abi, _, err := cpchain.ReadContract("../fixtures/contract/AirDrop.json")
	if err != nil {
		t.Fatal(err)
	}

	contractairdrop := client.Contract([]byte(Abi), airdropaddress)
	// tx, err := wallet.Transfer("0x2D770FC4B2E8F24292B08cc5fB15E6a69Fc0356a", int64(1))
	// t.Logf("Tx hash: %v", tx.Hash().Hex())

	tx, err := contractairdrop.Call(wallet, 41, "recharge", int64(20))
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

func Test3(t *testing.T) { //设置manager
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := client.LoadWallet(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}

	Abi, _, err := cpchain.ReadContract("../fixtures/contract/AirDrop.json")
	if err != nil {
		t.Fatal(err)
	}

	contractairdrop := client.Contract([]byte(Abi), airdropaddress)
	// tx, err := wallet.Transfer("0x2D770FC4B2E8F24292B08cc5fB15E6a69Fc0356a", int64(1))
	// t.Logf("Tx hash: %v", tx.Hash().Hex())

	tx, err := contractairdrop.Call(wallet, 41, "setManager", int64(0), "0xfd15c2932a60631222f7e6ffdde7bdab7237c2dc")
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

func Test4(t *testing.T) { //查看manager
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}

	Abi, _, err := cpchain.ReadContract("../fixtures/contract/AirDrop.json")
	if err != nil {
		t.Fatal(err)
	}

	contractairdrop := client.Contract([]byte(Abi), airdropaddress)

	tx, err := contractairdrop.View("isManager", "0xfd15c2932a60631222f7e6ffdde7bdab7237c2dc")
	fmt.Println(tx)
}

func Test5(t *testing.T) { //调用空投
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := client.LoadWallet(keystorePath, password)
	if err != nil {
		t.Fatal(err)
	}

	Abi, _, err := cpchain.ReadContract("../fixtures/contract/AirDrop.json")
	if err != nil {
		t.Fatal(err)
	}

	contractairdrop := client.Contract([]byte(Abi), airdropaddress)
	// tx, err := wallet.Transfer("0x2D770FC4B2E8F24292B08cc5fB15E6a69Fc0356a", int64(1))
	// t.Logf("Tx hash: %v", tx.Hash().Hex())

	tx, err := contractairdrop.Call(wallet, 41, "provideAirDrop", int64(0), "0xfd15c2932a60631222f7e6ffdde7bdab7237c2dc", "0x4f5625efef254760301d2766c6cc98f05722963e")
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
