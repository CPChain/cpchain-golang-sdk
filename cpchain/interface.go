package cpchain

import (
	"cpchain-golang-sdk/internal/fusion"
	"cpchain-golang-sdk/internal/fusion/contract"
	"math/big"
)

type Network struct {
	Name       string
	JsonRpcUrl string
	ChainId    uint
}

type EventsOptions struct {
	FromBlock uint64
	ToBlock   uint64
}

type WithEventsOptionsOption func(*EventsOptions)

func WithEventsOptionsFromBlock(block uint64) WithEventsOptionsOption {
	return func(flo *EventsOptions) {
		flo.FromBlock = block
	}
}

func WithEventsOptionsToBlock(block uint64) WithEventsOptionsOption {
	return func(flo *EventsOptions) {
		flo.ToBlock = block
	}
}

type Contract interface {
	Events(eventName string, event interface{}, options ...WithEventsOptionsOption) ([]*contract.Event, error)
}

type CPChain interface {
	// Get the current block number
	BlockNumber() (uint64, error)
	Block(number int) (*fusion.FullBlock, error)
	// get balance
	BalanceOf(address string) *big.Int
	Contract(abi []byte, address string) Contract
}
