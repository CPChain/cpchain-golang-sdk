package bind

import (
	"context"
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/cpcclient"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

// ContractCaller defines the methods needed to allow operating with contract on a read
// only basis.
type ContractCaller interface {
	// CodeAt returns the code of the given account. This is needed to differentiate
	// between contract internal errors and the local chain being out of sync.
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	// ContractCall executes an cpchain contract call with the specified data as the
	// input.
	CallContract(ctx context.Context, call cpcclient.CallMsg, blockNumber *big.Int) ([]byte, error)
}

// ContractTransactor defines the methods needed to allow operating with contract
// on a write only basis. Beside the transacting method, the remainder are helpers
// used when the user does not provide some needed values, but rather leaves it up
// to the transactor to decide.
type ContractTransactor interface {
	// PendingCodeAt returns the code of the given account in the pending state.
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	// PendingNonceAt retrieves the current pending nonce associated with an account.
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
	// execution of a transaction.
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	// EstimateGas tries to estimate the gas needed to execute a specific
	// transaction based on the current pending state of the backend blockchain.
	// There is no guarantee that this is the true gas limit requirement as other
	// transactions may be added or removed by miners, but it should provide a basis
	// for setting a reasonable default.
	EstimateGas(ctx context.Context, call cpcclient.CallMsg) (gas uint64, err error)
	// SendTransaction injects the transaction into the pending pool for execution.
	SendTransaction(ctx context.Context, tx *types.Transaction) error
}

// ContractFilterer defines the methods needed to access log events using one-off
// queries or continuous event subscriptions.
type ContractFilterer interface {
	// FilterLogs executes a log filter operation, blocking during execution and
	// returning all the results in one batch.
	//
	// TODO(karalabe): Deprecate when the subscription one can return past data too.
	FilterLogs(ctx context.Context, query types.FilterQuery) ([]types.Log, error)
}

// ContractBackend defines the methods needed to work with contracts on a read-write basis.
type ContractBackend interface {
	ContractFilterer
	ContractCaller
	ContractTransactor
}
