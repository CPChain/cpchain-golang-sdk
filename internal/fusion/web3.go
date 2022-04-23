package fusion

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"cpchain-golang-sdk/internal/fusion/hexutil"

	"github.com/zgljl2012/slog"
)

type web3 struct {
	provider Provider
}

func NewWeb3(provider Provider) (Web3, error) {
	return &web3{
		provider: provider,
	}, nil
}

func handleRawBlockWithFullTxs(rawBlock *RawBlockWithFullTxs) *FullBlock {
	var block FullBlock
	block.Hash = rawBlock.Hash
	block.GasLimit, _ = hexutil.DecodeUint64(string(rawBlock.GasLimit))
	block.GasUsed, _ = hexutil.DecodeUint64(string(rawBlock.GasUsed))
	block.ExtraData = rawBlock.Hash
	block.LogsBloom = rawBlock.LogsBloom
	block.Miner = string(rawBlock.Miner)
	block.Number, _ = hexutil.DecodeUint64(string(rawBlock.Number))
	block.ParentHash = rawBlock.ParentHash
	block.ReceiptsRoot = rawBlock.ReceiptsRoot
	block.Size, _ = hexutil.DecodeUint64(string(rawBlock.Size))
	block.StateRoot = rawBlock.StateRoot
	block.Timestamp, _ = hexutil.DecodeUint64(string(rawBlock.Timestamp))
	block.TransactionsRoot = rawBlock.TransactionsRoot
	// handle tx
	for _, rawTx := range rawBlock.Transactions {
		var tx Tx
		tx.BlockHash = rawTx.BlockHash
		tx.BlockNumber, _ = hexutil.DecodeUint64(string(rawTx.BlockNumber))
		tx.From = rawTx.From
		tx.Gas, _ = hexutil.DecodeUint64(string(rawTx.Gas))
		tx.GasPrice, _ = hexutil.DecodeUint64(string(rawTx.GasPrice))
		tx.Hash = rawTx.Hash
		tx.Input = rawTx.Input
		tx.Nonce, _ = hexutil.DecodeUint64(string(rawTx.Nonce))
		tx.R = rawTx.R
		tx.S = rawTx.S
		tx.V = rawTx.V
		tx.To = rawTx.To
		tx.TransactionIndex, _ = hexutil.DecodeUint64(string(rawTx.TransactionIndex))
		tx.Type, _ = hexutil.DecodeUint64(string(rawTx.Type))
		value := new(big.Int)
		s := rawTx.Value[2:]
		value.SetString(string(s), 16)
		tx.Value = value
		block.Transactions = append(block.Transactions, tx)
	}
	return &block
}

func handleNumber(number interface{}) (interface{}, error) {
	_num := number
	if reflect.TypeOf(number).Kind() == reflect.String {
		if number.(string) != "latest" {
			return nil, fmt.Errorf("unknown type: %s", number)
		}
	} else if reflect.TypeOf(number).Kind() == reflect.Uint64 {
		_num = hexutil.EncodeUint64(number.(uint64))
	} else if reflect.TypeOf(number).Kind() == reflect.Int {
		_num = hexutil.EncodeUint64(uint64(number.(int)))
	}
	return _num, nil
}

func (w *web3) GetBlock(number interface{}) (interface{}, error) {
	_num, err := handleNumber(number)
	if err != nil {
		return nil, fmt.Errorf("get block failed: %s", err)
	}
	data, err := w.provider.MakeRequest("eth_getBlockByNumber", []interface{}{_num, true})
	if err != nil {
		return nil, fmt.Errorf("get block failed: %s", err)
	}
	var rawBlock RawBlockWithFullTxs
	if err := json.Unmarshal(data, &rawBlock); err != nil {
		return nil, fmt.Errorf("unmarshal raw-block failed: %s", err)
	}
	return handleRawBlockWithFullTxs(&rawBlock), nil
}

func (w *web3) GetBlockByNumber(number interface{}, fullTx bool) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (w *web3) GetBalanceAt(address string, number interface{}) (*big.Int, error) {
	_num, err := handleNumber(number)
	if err != nil {
		return nil, fmt.Errorf("get balance failed: %s", err)
	}
	data, err := w.provider.MakeRequest("eth_getBalance", []interface{}{address, _num})
	if err != nil {
		return nil, fmt.Errorf("get balance failed: %s", err)
	}
	balance := string(data)
	balance = balance[3 : len(balance)-1]
	slog.Info("test", "1", string(data), "2", balance)

	r := big.NewInt(0)
	r, _ = r.SetString(balance, 16)
	return r, nil
}

func (w *web3) GetBalance(address string) *big.Int {
	if r, err := w.GetBalanceAt(address, "latest"); err != nil {
		slog.Error(err)
	} else {
		return r
	}
	return big.NewInt(0)
}
