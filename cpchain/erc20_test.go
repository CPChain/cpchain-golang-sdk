package cpchain_test

import (
	"math/big"
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/cpchain/modules/token"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
)

// 在测试链上部署一个 ERC20 合约，进行事件测试
const (
	// 合约地址
	ContractAddress = "0x3e9ee62921AE6af4341AE0E923c114Fae55fdC38"
	// 合约名称
	ContractName = "Meta Token"
	// 合约简称
	ContractSymbol = "MT"
	// 合约单位
	ContractUnit = "MT"
	// 合约精度
	ContractDecimals = 18
	// 合约发行总量
	ContractTotalSupply = 10000
	// 合约接口文件路径
	ContractABIPath = "../fixtures/MetaToken.json"
)

func TestMetaToken(t *testing.T) {
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	contract := token.NewERC20Contract(client, ContractAddress)
	events, err := contract.Events(token.TRANSFER_EVENT_NAME, token.TransferEvent{})
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events {
		transferEvent := event.Data.(*token.TransferEvent)
		t.Logf("TransferEvent: %+v", transferEvent)
		t.Logf("Blocknumber: %d, From: %v To: %v, Value: %d",
			event.BlockNumber, transferEvent.From.Hex(), transferEvent.To.Hex(), big.NewInt(0).Div(transferEvent.Value, big.NewInt(1e18)))
	}
}

func TestMetaTokenWithMap(t *testing.T) {
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	contract := token.NewERC20Contract(client, ContractAddress)
	events, err := contract.Events(token.TRANSFER_EVENT_NAME, map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events {
		transferEvent := event.Data.(map[string]interface{})
		t.Logf("TransferEvent: %+v", transferEvent)
		val := transferEvent["value"].(big.Int)
		t.Logf("Blocknumber: %d, From: %v To: %v, Value: %d",
			event.BlockNumber, transferEvent["from"].(common.Address).Hex(), transferEvent["to"].(common.Address).Hex(), big.NewInt(0).Div(&val, big.NewInt(1e18)))
	}
}

func TestHandGame(t *testing.T) {
	// Handgame 中包含一个 uint64 的 indexed field，所以测试一下
	var (
		abi             = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"limit","type":"uint256"}],"name":"SetMaxLimit","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"limit","type":"uint256"}],"name":"SetTimeoutLimit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"gameId","type":"uint64"},{"indexed":false,"name":"starter","type":"address"},{"indexed":false,"name":"card","type":"uint256"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"threshold","type":"uint256"}],"name":"GameStarted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"group_id","type":"uint256"},{"indexed":true,"name":"gameId","type":"uint64"},{"indexed":false,"name":"starter","type":"address"},{"indexed":false,"name":"message","type":"string"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"threshold","type":"uint256"}],"name":"CreateGroupHandGame","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"gameId","type":"uint64"}],"name":"GameCancelled","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"gameId","type":"uint64"},{"indexed":false,"name":"player","type":"address"},{"indexed":false,"name":"card","type":"uint256"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"GameLocked","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"gameId","type":"uint64"},{"indexed":false,"name":"player","type":"address"},{"indexed":false,"name":"key","type":"uint256"},{"indexed":false,"name":"content","type":"uint256"}],"name":"CardOpened","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"gameId","type":"uint64"},{"indexed":false,"name":"result","type":"int8"}],"name":"GameFinished","type":"event"},{"constant":false,"inputs":[{"name":"limit","type":"uint256"}],"name":"setMaxLimit","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"limit","type":"uint256"}],"name":"setTimeoutLimit","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
		contractAddress = "0xd9Fff84402f066157b723310fa755B2209346910"
	)

	client, err := cpchain.NewCPChain(cpchain.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	contract := client.Contract([]byte(abi), contractAddress)
	events, err := contract.Events("CardOpened", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(events))
	for _, e := range events {
		cardOpened := e.Data.(map[string]interface{})
		t.Log(cardOpened)
	}
}
