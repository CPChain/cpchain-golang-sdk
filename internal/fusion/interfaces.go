package fusion

import (
	"math/big"
)

type Provider interface {
	MakeRequest(method string, args []interface{}) ([]byte, error)
}

type Web3 interface {
	GetBlock(number interface{}) (interface{}, error) // number is a digit or 'latest'
	GetBlockByNumber(number interface{}, fullTx bool) (interface{}, error)
	// number: "latest"/or number
	GetBalanceAt(address string, number interface{}) (*big.Int, error)
	GetBalance(address string) *big.Int
}
