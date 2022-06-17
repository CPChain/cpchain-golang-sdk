package cpchain

import (
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/contract"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

type Event = contract.Event

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
	// TODO 如果事件非常多，如10000条事件，是否会分批获取？
	Events(eventName string, event interface{}, options ...WithEventsOptionsOption) ([]*contract.Event, error)
}

type CPChain interface {
	// Get the current block number
	BlockNumber() (uint64, error)
	Block(number int) (*fusion.FullBlock, error)
	// get balance
	BalanceOf(address string) *big.Int
	Contract(abi []byte, address string) Contract
	//load a wallet by keystore path
	LoadWallet(path string) Wallet //TODO 是否要加error
	//create a wallet by dirpath and password
	CreateWallet(path string, password string) (*Account, error) //返回值或需更改
	//deploy contract to chain
	DeployContract() (common.Address, error)
}

// TODO simulate chain

type Wallet interface {
	Addr() common.Address

	GetKey(password string) (*Key, error)

	SignTxWithPassword(password string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}
