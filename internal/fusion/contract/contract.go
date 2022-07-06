package contract

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/CPChain/cpchain-golang-sdk/internal/cpcclient"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind/backends"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

var (
	// ErrNoCode is returned by call and transact operations for which the requested
	// recipient contract to operate on does not exist in the state db or does not
	// have any code associated with it (i.e. suicided).
	ErrNoCode = errors.New("no contract code at given address")

	// This error is raised when attempting to perform a pending state action
	// on a backend that doesn't implement PendingContractCaller.
	ErrNoPendingState = errors.New("backend does not support pending state")

	// This error is returned by WaitDeployed if contract creation leaves an
	// empty contract behind.
	ErrNoCodeAfterDeploy = errors.New("no contract code after deployment")
)

type FilterLogsOptions struct {
	FromBlock uint64
	ToBlock   uint64
}

type WithFilterLogsOption func(*FilterLogsOptions)

func WithFilterLogsFromBlock(block uint64) WithFilterLogsOption {
	return func(flo *FilterLogsOptions) {
		flo.FromBlock = block
	}
}

func WithFilterLogsEndBlock(block uint64) WithFilterLogsOption {
	return func(flo *FilterLogsOptions) {
		flo.ToBlock = block
	}
}

type Event struct {
	types.Log
	Name string
	Data interface{}
}

type Contract interface {
	FilterLogs(eventName string, event interface{}, options ...WithFilterLogsOption) ([]*Event, error) // event parameter is a event struct, e.g. CreateProduct{}
	// Depoly(opts *bind.TransactOpts, bytecode []byte, chainId uint, params ...interface{}) (common.Address, *types.Transaction, *contract, error)
	Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error

	Transact(opts *bind.TransactOpts, chainId uint, method string, params ...interface{}) (*types.Transaction, error)
}

type contract struct { //NOTE:like boundcontract
	abi     abi.ABI
	address common.Address
	backend bind.ContractBackend
}

// Contract
func NewContractWithProvider(abi []byte, address common.Address, provider fusion.Provider) (Contract, error) {
	backend := backends.NewClientBackend(provider)
	return NewContract(abi, address, backend)
}

func NewContractWithUrl(abi []byte, address common.Address, url string) (Contract, error) {
	backend, err := cpcclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return NewContract(abi, address, backend)
}

func NewContract(abiData []byte, address common.Address, backend bind.ContractBackend) (Contract, error) {
	instance, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		return nil, fmt.Errorf("new contract failed: %v", err)
	}
	return &contract{
		abi:     instance,
		address: address,
		backend: backend,
	}, nil
}

func (c *contract) FilterLogs(eventName string, event interface{}, options ...WithFilterLogsOption) ([]*Event, error) {
	var opts = FilterLogsOptions{}
	for _, op := range options {
		op(&opts)
	}
	query := [][]interface{}{{c.abi.Events[eventName].Id()}}

	topics, err := bind.MakeTopics(query...)
	if err != nil {
		return nil, fmt.Errorf("filter logs failed: %v", err)
	}
	filterQuery := types.FilterQuery{
		Addresses: []common.Address{c.address},
		Topics:    topics,
		FromBlock: new(big.Int).SetUint64(opts.FromBlock),
	}
	logs, err := c.backend.FilterLogs(context.Background(), filterQuery)
	if err != nil {
		return nil, fmt.Errorf("filter logs failed: %v", err)
	}
	var events []*Event
	vt := reflect.TypeOf(event)
	eventIsMap := vt.Kind() == reflect.Map
	for _, l := range logs {
		var v reflect.Value
		if vt.Kind() != reflect.Map {
			v = reflect.New(vt)
			if err := c.abi.Unpack(v.Interface(), eventName, l.Data); err != nil {
				return nil, fmt.Errorf("unpack event failed (not map): %v", err)
			}
		} else {
			v = reflect.MakeMap(vt)
			tmp := map[string]interface{}{}
			if err := c.abi.Unpack(&tmp, eventName, l.Data); err == nil {
				for k_, v_ := range tmp {
					vv := reflect.ValueOf(v_)
					vv = reflect.Indirect(vv)
					v.SetMapIndex(reflect.ValueOf(k_), vv)
				}
			} else {
				// 允许当事件中无任何索引字段时，返回空 map
				if err.Error() != "abi: unmarshalling empty output" {
					return nil, fmt.Errorf("unpack event failed (map): %v", err)
				}
			}
		}
		// 处理 Indexed 字段
		if len(l.Topics) > 1 {
			if err := c.abi.Events[eventName].Inputs.ForEach(func(i int, input abi.Argument) error {
				if input.Indexed {
					val := l.Topics[i+1]
					// TODO support other types
					if input.Type.String() == "address" {
						if !eventIsMap {
							v.Elem().Field(i).Set(reflect.ValueOf(common.HexToAddress(val.Hex())))
						} else {
							v.SetMapIndex(reflect.ValueOf(input.Name), reflect.ValueOf(common.HexToAddress(val.Hex())))
						}
					} else if input.Type.String() == "uint64" {
						if !eventIsMap {
							v.Elem().Field(i).Set(reflect.ValueOf(new(big.Int).SetUint64(val.Big().Uint64())))
						} else {
							v.SetMapIndex(reflect.ValueOf(input.Name), reflect.ValueOf(new(big.Int).SetUint64(val.Big().Uint64())))
						}
					}
				}
				return nil
			}); err != nil {
				return nil, fmt.Errorf("handle indexed fields failed: %v", err)
			}
		}
		events = append(events, &Event{
			Name: eventName,
			Log:  l,
			Data: v.Interface(),
		})
	}
	return events, nil
}

// view
func (c *contract) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	// Don't crash on a lazy user
	if opts == nil {
		opts = new(bind.CallOpts)
	}
	// Pack the input, call and unpack the results
	input, err := c.abi.Pack(method, params...)
	if err != nil {
		return err
	}
	var (
		msg    = cpcclient.CallMsg{From: opts.From, To: &c.address, Data: input}
		ctx    = ensureContext(opts.Context)
		code   []byte
		output []byte
	)
	if opts.Pending {
		pb, ok := c.backend.(bind.PendingContractCaller)
		if !ok {
			return ErrNoPendingState
		}
		output, err = pb.PendingCallContract(ctx, msg)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			// NOTE: it may cause some edge case where the pending block doesn't add to chain and the tx which depends on it will eventually fail
			if code, err = pb.PendingCodeAt(ctx, c.address); err != nil {
				return err
			} else if len(code) == 0 {
				return ErrNoCode
			}
		}
	} else {
		output, err = c.backend.CallContract(ctx, msg, nil)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = c.backend.CodeAt(ctx, c.address, nil); err != nil {
				return err
			} else if len(code) == 0 {
				return ErrNoCode
			}
		}
	}
	if err != nil {
		return err
	}
	return c.abi.Unpack(result, method, output)
}

// Transact invokes the (paid) contract method with params as input values.
func (c *contract) Transact(opts *bind.TransactOpts, chainId uint, method string, params ...interface{}) (*types.Transaction, error) {
	// Otherwise pack up the parameters and invoke the contract
	input, err := c.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}
	return c.transact(chainId, opts, &c.address, input)
}

//TODO chainID
// transact executes an actual transaction invocation, first deriving any missing
// authorization fields, and then scheduling the transaction for execution.
func (c *contract) transact(chainId uint, opts *bind.TransactOpts, contract *common.Address, input []byte) (*types.Transaction, error) {
	var err error
	// Ensure a valid value field and resolve the account nonce
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	var nonce uint64
	if opts.Nonce == nil {
		nonce, err = c.backend.PendingNonceAt(ensureContext(opts.Context), opts.From)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
		}
	} else {
		nonce = opts.Nonce.Uint64()
	}
	// Figure out the gas allowance and gas price values
	gasPrice := opts.GasPrice
	if gasPrice == nil {
		gasPrice, err = c.backend.SuggestGasPrice(ensureContext(opts.Context))
		if err != nil {
			return nil, fmt.Errorf("failed to suggest gas price: %v", err)
		}
	}
	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		// Gas estimation cannot succeed without code for method invocations
		if contract != nil {
			if code, err := c.backend.PendingCodeAt(ensureContext(opts.Context), c.address); err != nil {
				return nil, err
			} else if len(code) == 0 {
				return nil, ErrNoCode
			}
		}
		// If the contract surely has code (or code is not needed), estimate the transaction
		msg := cpcclient.CallMsg{From: opts.From, To: contract, Value: value, Data: input}
		gasLimit, err = c.backend.EstimateGas(ensureContext(opts.Context), msg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas needed: %v", err)
		}
	}
	// Create the transaction, sign it and schedule it for execution
	var rawTx *types.Transaction
	if contract == nil {
		rawTx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, input)
	} else {
		rawTx = types.NewTransaction(nonce, c.address, value, gasLimit, gasPrice, input)
	}
	if opts.Signer == nil {
		return nil, errors.New("no signer to authorize the transaction with")
	}

	ChainID := big.NewInt(int64(chainId)) //TODO

	// signedTx, err := opts.Signer(types.NewCep1Signer(configs.ChainConfigInfo().ChainID), opts.From, rawTx)
	signedTx, err := opts.Signer(types.NewCep1Signer(ChainID), opts.From, rawTx)

	if err != nil {
		return nil, err
	}
	if err := c.backend.SendTransaction(ensureContext(opts.Context), signedTx); err != nil {
		return nil, err
	}
	return signedTx, nil
}

// ensureContext is a helper method to ensure a context is not nil, even if the
// user specified it as such.
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}
