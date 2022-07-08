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

// type clientBackend struct {
// 	provider fusion.Provider
// }
type clientBackend struct {
	provider fusion.Provider
	client   cpcclient.Client
}

// func NewClientBackend(provider fusion.Provider) bind.ContractBackend {
// 	return &clientBackend{
// 		provider: provider,
// 	}
// }

func NewClientBackend(provider fusion.Provider, endpoint string) bind.ContractBackend {
	client, err := cpcclient.Dial(endpoint)
	if err != nil {
		return nil
	}
	return &clientBackend{
		provider: provider,
		client:   *client,
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
	result, err := c.client.PendingNonceAt(ctx, account)
	return result, err
}

func (c *clientBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	result, err := c.client.PendingCodeAt(ctx, account)
	return result, err
}

func (c *clientBackend) PendingCallContract(ctx context.Context, msg cpcclient.CallMsg) ([]byte, error) {
	result, err := c.client.PendingCallContract(ctx, msg)
	return result, err
}

func (c *clientBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	result, err := c.client.SuggestGasPrice(ctx)
	return result, err
}

func (c *clientBackend) EstimateGas(ctx context.Context, msg cpcclient.CallMsg) (uint64, error) {
	result, err := c.client.EstimateGas(ctx, msg)
	return result, err
}

func (c *clientBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	err := c.client.SendTransaction(ctx, tx)
	return err
}

func (c *clientBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	result, err := c.client.CodeAt(ctx, contract, blockNumber)
	return result, err
}

// ContractCall executes an cpchain contract call with the specified data as the
// input.
func (c *clientBackend) CallContract(ctx context.Context, call cpcclient.CallMsg, blockNumber *big.Int) ([]byte, error) {
	result, err := c.client.CallContract(ctx, call, blockNumber)
	return result, err
}

func (c *clientBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	result, err := c.client.TransactionReceipt(ctx, txHash)
	return result, err
}
