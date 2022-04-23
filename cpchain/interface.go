package cpchain

import (
	"cpchain-golang-sdk/internal/fusion"
	"math/big"
)

type Network struct {
	Name       string
	JsonRpcUrl string
	ChainId    uint
}

type CPChain interface {
	// Get the current block number
	BlockNumber() (uint64, error)
	Block(number int) (*fusion.FullBlock, error)
	// GetBalanceAt(address string, number interface{}) (*big.Int, error)
	BalanceOf(address string) *big.Int
}
