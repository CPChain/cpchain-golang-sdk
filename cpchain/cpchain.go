package cpchain

import (
	"cpchain-golang-sdk/internal/fusion"
	"fmt"
)

type cpchain struct {
	network  Network
	provider fusion.Provider
	web3     fusion.Web3
}

func NewCPChain(network Network) (CPChain, error) {
	provider, err := fusion.NewHttpProvider(network.JsonRpcUrl)
	if err != nil {
		return nil, err
	}
	web3, err := fusion.NewWeb3(provider)
	if err != nil {
		return nil, err
	}
	return &cpchain{
		provider: provider,
		network:  network,
		web3:     web3,
	}, nil
}

func (c *cpchain) GetBlockNumber() (uint64, error) {
	block, err := c.web3.GetBlock("latest")
	if err != nil {
		return 0, fmt.Errorf("get block number failed: %v", err)
	}
	return block.(*fusion.FullBlock).Number, nil
}

// func (c *cpchain) GetBlock(number interface{}) (interface{}, error) {

// }

// func (c *cpchain) GetBlockByNumber(number interface{}, fullTx bool) (interface{}, error) {

// }

// func (c *cpchain) GetBalanceAt(address string, number interface{}) (*big.Int, error) {

// }

// func (c *cpchain) GetBalance(address string) *big.Int {

// }
