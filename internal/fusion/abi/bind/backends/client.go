package backends

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/cpcclient"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common/hexutil"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

type clientBackend struct {
	provider fusion.Provider
}

func NewClientBackend(provider fusion.Provider) bind.ContractBackend {
	return &clientBackend{
		provider: provider,
	}
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

func toFilterArg(q types.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   q.Addresses,
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

func rawlogs2logs(rawlogs []rawlog) []types.Log {
	if rawlogs == nil {
		return nil
	}
	var r []types.Log
	for _, l := range rawlogs {
		var topics []common.Hash
		for _, t := range l.Topics {
			topics = append(topics, common.HexToHash(t))
		}
		bn, _ := hexutil.DecodeUint64(l.BlockNumber)
		tindex, _ := hexutil.DecodeUint64(l.TransactionIndex)
		index, _ := hexutil.DecodeUint64(l.LogIndex)
		r = append(r, types.Log{
			Address:     common.HexToAddress(l.Address),
			Topics:      topics,
			Data:        common.FromHex(l.Data),
			BlockNumber: bn,
			TxHash:      common.HexToHash(l.TransactionHash),
			TxIndex:     uint(tindex),
			BlockHash:   common.HexToHash(l.BlockHash),
			Index:       uint(index),
			Removed:     l.Removed,
		})
	}
	return r
}

func (c *clientBackend) FilterLogs(ctx context.Context, query types.FilterQuery) ([]types.Log, error) {
	args := toFilterArg(query)
	results, err := c.provider.MakeRequest("eth_getLogs", []interface{}{args})
	if err != nil {
		return nil, fmt.Errorf("Provider make request failed: %v", err)
	}
	var data []rawlog
	if err := json.Unmarshal(results, &data); err != nil {
		return nil, fmt.Errorf("Unmarshal failed: %v", err)
	}
	return rawlogs2logs(data), nil
}

// TODO 以后将clientbackend 替换掉cpcclient
func (c *clientBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {

	return uint64(0), nil
}

func (c *clientBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {

	return nil, nil
}

func (c *clientBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return nil, nil
}

func (c *clientBackend) EstimateGas(ctx context.Context, msg cpcclient.CallMsg) (uint64, error) {
	return uint64(0), nil
}

func (c *clientBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}
