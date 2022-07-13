package fusion

import (
	"math/big"
	"context"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
)

type Provider interface {
	MakeRequest(method string, args []interface{}) ([]byte, error)
}

type Web3 interface {
	GetBlock(number interface{}) (interface{}, error) // number is a digit or 'latest'
	GetBlockByNumber(number interface{}, fullTx bool) (interface{}, error)
	// number: "latest"/or number
	GetBalanceAt(address string, number interface{}) (*big.Int, error)
	GetBalance(address string) *big.Int
}

type ChainStateReader interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
}