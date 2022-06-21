package contract

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

// The difference between NewBoundContract and NewContract is that the return type is different
func NewBoundContract(abiData string, address common.Address, backend bind.ContractBackend) (*contract, error) {
	instance, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		return nil, fmt.Errorf("new bound contract failed: %v", err)
	}
	return &contract{
		abi:     instance,
		address: address,
		backend: backend,
	}, nil
}

func DeployContract(abiData string, opts *bind.TransactOpts, bytecode []byte, backend bind.ContractBackend, params ...interface{}) (common.Address, *types.Transaction, *contract, error) {
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

// from chain/tools/smartcontract/depoly/utils
func NewTransactor(privateKey *ecdsa.PrivateKey, nonce *big.Int) *bind.TransactOpts {
	auth := bind.NewKeyedTransactor(privateKey)
	if nonce.Cmp(big.NewInt(-1)) > 0 {
		auth.Nonce = nonce
	}
	return auth
}
