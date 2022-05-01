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
