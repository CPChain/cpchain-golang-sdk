package cpchain

import (
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind"
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
	//Generate an account based on the password and store its keystore to the path
	CreateAccount(path string, password string) (*Account, error)
	//return backend
	Backend() (bind.ContractBackend, error)
}

// TODO simulate chain

type Wallet interface {
	// 返回钱包地址
	Addr() common.Address

	// 获取密钥
	GetKey(password string) (*Key, error)

	// 给交易签名
	SignTxWithPassword(password string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	// 交易
	Transfer(password string, targetAddr string, value int64) error //TODO 返回值

	// 通过文件部署合约
	DeployContractByFile(path string, password string) error

	// 通过abi和bin部署合约
	DeployContract(abi string, bin string, password string) error
}
