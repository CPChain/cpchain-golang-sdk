package contract

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"reflect"

	"cpchain-golang-sdk/internal/fusion"
	"cpchain-golang-sdk/internal/fusion/abi"
	"cpchain-golang-sdk/internal/fusion/abi/bind"
	"cpchain-golang-sdk/internal/fusion/abi/bind/backends"
	"cpchain-golang-sdk/internal/fusion/common"
	"cpchain-golang-sdk/internal/fusion/types"

	"github.com/zgljl2012/slog"
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
	Data interface{}
}

type Contract interface {
	FilterLogs(eventName string, event interface{}, options ...WithFilterLogsOption) ([]*Event, error) // event parameter is a event struct, e.g. CreateProduct{}
}

type contract struct {
	abi     abi.ABI
	address common.Address
	backend bind.ContractBackend
}

// Contract
func NewContractWithProvider(abi []byte, address common.Address, provider fusion.Provider) (Contract, error) {
	backend := backends.NewClientBackend(provider)
	return NewContract(abi, address, backend)
}

func NewContract(abiData []byte, address common.Address, backend bind.ContractBackend) (Contract, error) {
	instance, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		return nil, fmt.Errorf("New contract failed: %v", err)
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
		return nil, fmt.Errorf("Filter logs failed: %v", err)
	}
	filterQuery := types.FilterQuery{
		Addresses: []common.Address{c.address},
		Topics:    topics,
		FromBlock: new(big.Int).SetUint64(opts.FromBlock),
	}
	logs, err := c.backend.FilterLogs(context.Background(), filterQuery)
	if err != nil {
		return nil, fmt.Errorf("Filter logs failed: %v", err)
	}
	var events []*Event
	vt := reflect.TypeOf(event)
	for _, l := range logs {
		v := reflect.New(vt)
		if err := c.abi.Unpack(v.Interface(), eventName, l.Data); err != nil {
			slog.Error(err)
		}
		events = append(events, &Event{
			Log:  l,
			Data: v.Interface(),
		})
	}
	return events, nil
}
