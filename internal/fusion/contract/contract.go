package contract

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"reflect"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind/backends"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
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
