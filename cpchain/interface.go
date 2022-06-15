package cpchain

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/contract"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto/ecies"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
	"github.com/pborman/uuid"
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
	LoadWallet(path string) Wallet
	//create a wallet by dirpath and password
	CreateWallet(path string, password string) Wallet //返回值或需更改
}

// TODO simulate chain

type Wallet interface {
	Addr() common.Address

	GetKey(password string) (*Key, error)

	SignTxWithPassword(password string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}

type Key struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey *ecdsa.PrivateKey

	// we cache ecies.PrivateKey used to decrypt data that encrypted by ecies.PublicKey from remote node
	EciesPrivateKey *ecies.PrivateKey
}
