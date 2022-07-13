package cpchain

import (
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/contract"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
	"github.com/CPChain/cpchain-golang-sdk/internal/keystore"
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

	// Call function, need to send transation
	Call(w Wallet, chainId uint, method string, params ...string) (*types.Transaction, error)

	// View
	View(method string, params ...string) (interface{}, error)
}

type CPChain interface {
	// Get the current block number
	BlockNumber() (uint64, error)
	//
	Block(number int) (*fusion.FullBlock, error)
	// Get nonce of account
	NonceOf(address string) ([]byte, error)
	// Get gas price
	GasPrice() (*big.Int, error)
	// Get balance
	BalanceOf(address string) *big.Int
	// New Contract instance
	Contract(abi []byte, address string) Contract
	// Load a wallet by keystore path
	LoadWallet(path string, password string) (Wallet, error) // TODO 要加error
	// Generate an account based on the password and store its keystore to the path
	CreateAccount(path string, password string) (*Account, error)
	// 通过文件部署合约
	DeployContractByFile(path string, w Wallet) (common.Address, *types.Transaction, error)
	// 通过abi和bin部署合约
	DeployContract(abi string, bin string, w Wallet) (common.Address, *types.Transaction, error)
	// 通过签名过的交易获取结果
	ReceiptByTx(signedTx *types.Transaction) (*types.Receipt, error)
}

// TODO simulate chain

type Wallet interface {
	// 返回钱包地址
	Addr() common.Address

	// 获取密钥
	Key() *keystore.Key

	// sign transaction
	SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	// 交易
	Transfer(targetAddr string, value int64) (*types.Transaction, error) //TODO 返回值
}
