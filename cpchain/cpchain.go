package cpchain

import (
	"cpchain-golang-sdk/internal/fusion"
	"fmt"
	"math/big"
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

func (c *cpchain) BlockNumber() (uint64, error) {
	block, err := c.web3.GetBlock("latest")
	if err != nil {
		return 0, fmt.Errorf("get block number failed: %v", err)
	}
	return block.(*fusion.FullBlock).Number, nil
}

func (c *cpchain) Block(number int) (*fusion.FullBlock, error) {
	block, err := c.web3.GetBlock(uint64(number))
	if err != nil {
		return nil, fmt.Errorf("get block failed: %v", err)
	}
	return block.(*fusion.FullBlock), nil
}

func (c *cpchain) BalanceOf(address string) *big.Int {
	balance, err := c.web3.GetBalanceAt(address, "latest")
	if err != nil {
		return big.NewInt(0)
	}
	return balance
}

func WeiToCpc(wei *big.Int) *big.Int {
	return wei.Div(wei, big.NewInt(1e18))
}

// func (c *cpchain) GetBlockByNumber(number interface{}, fullTx bool) (interface{}, error) {

// }

// func (c *cpchain) GetBalanceAt(address string, number interface{}) (*big.Int, error) {

// }

// func (c *cpchain) GetBalance(address string) *big.Int {

// }
