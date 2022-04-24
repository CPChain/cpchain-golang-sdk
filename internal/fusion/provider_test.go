package fusion_test

import (
	"encoding/json"
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/hexutil"
)

func TestProvider(t *testing.T) {
	URL := "https://civilian.cpchain.io"
	provider, err := fusion.NewHttpProvider(URL)
	if err != nil {
		t.Fatal(err)
	}
	data, err := provider.MakeRequest("eth_getBlockByNumber", []interface{}{hexutil.EncodeUint64(6400839), false})
	if err != nil {
		t.Fatal(err)
	}
	var rawBlock fusion.RawBlock
	if err := json.Unmarshal(data, &rawBlock); err != nil {
		t.Fatal(err)
	}
	t.Log("dpor", rawBlock.Dpor)
	t.Log("extraData", rawBlock.ExtraData)
	t.Log("gasLimit", rawBlock.GasLimit)
	t.Log("gasUsed", rawBlock.GasUsed)
	t.Log("hash", rawBlock.Hash)
	t.Log("logsBloom", rawBlock.LogsBloom)
	t.Log("miner", rawBlock.Miner)
	t.Log("number", rawBlock.Number)
	t.Log("parentHash", rawBlock.ParentHash)
	t.Log("receiptsRoot", rawBlock.ReceiptsRoot)
	t.Log("size", rawBlock.Size)
	t.Log("stateRoot", rawBlock.StateRoot)
	t.Log("timestamp", rawBlock.Timestamp)
	t.Log("transactions", rawBlock.Transactions)
	t.Log("transactionsRoot", rawBlock.TransactionsRoot)
}

func TestGetBlocksWithFullTx(t *testing.T) {
	URL := "https://civilian.cpchain.io"
	provider, err := fusion.NewHttpProvider(URL)
	if err != nil {
		t.Fatal(err)
	}
	data, err := provider.MakeRequest("eth_getBlockByNumber", []interface{}{hexutil.EncodeUint64(6400839), true})
	if err != nil {
		t.Fatal(err)
	}
	var rawBlock fusion.RawBlockWithFullTxs
	if err := json.Unmarshal(data, &rawBlock); err != nil {
		t.Fatal(err)
	}
	var tx = rawBlock.Transactions[3]
	t.Log("blockHash", tx.BlockHash)
	t.Log("blockNumber", tx.BlockNumber)
	t.Log("from", tx.From)
	t.Log("gas", tx.Gas)
	t.Log("gasPrice", tx.GasPrice)
	t.Log("hash", tx.Hash)
	t.Log("input", tx.Input)
	t.Log("nonce", tx.Nonce)
	t.Log("r", tx.R)
	t.Log("s", tx.S)
	t.Log("v", tx.V)
	t.Log("to", tx.To)
	t.Log("transactionIndex", tx.TransactionIndex)
	t.Log("type", tx.Type)
	t.Log("v", tx.V)
	t.Log("value", tx.Value)
}
