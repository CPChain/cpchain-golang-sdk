package contract

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

// SignerFn is a signer function callback when a contract requires a method to
// sign the transaction before submission.

type SignerFn func(types.Signer, common.Address, *types.Transaction) (*types.Transaction, error)

// TransactOpts is the collection of authorization data required to create a
// valid Ethereum transaction.

type TransactOpts struct {
	From   common.Address // Ethereum account to send the transaction from
	Nonce  *big.Int       // Nonce to use for the transaction execution (nil = use pending state)
	Signer SignerFn       // Method to use for signing the transaction (mandatory)

	Value    *big.Int // Funds to transfer along along the transaction (nil = 0 = no funds)
	GasPrice *big.Int // Gas price to use for the transaction execution (nil = gas price oracle)
	GasLimit uint64   // Gas limit to set for the transaction execution (0 = estimate)

	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// The difference between NewBoundContract and NewContract is that the return type is different
func NewBoundContract(abiData []byte, address common.Address, backend bind.ContractBackend) (*contract, error) {
	instance, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		return nil, fmt.Errorf("new bound contract failed: %v", err)
	}
	return &contract{
		abi:     instance,
		address: address,
		backend: backend,
	}, nil
}

func DeployContract(abiData []byte, opts *TransactOpts, bytecode []byte, backend bind.ContractBackend, params ...interface{}) (common.Address, *types.Transaction, *contract, error) {
	c, err := NewBoundContract(abiData, common.Address{}, backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	input, err := c.abi.Pack("", params...)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	tx, err := c.transact(opts, nil, append(bytecode, input...))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	c.address = crypto.CreateAddress(opts.From, tx.Nonce())
	return c.address, tx, c, nil
}
