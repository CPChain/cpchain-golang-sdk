package cpchain

type Network struct {
	Name       string
	JsonRpcUrl string
	ChainId    uint
}

type CPChain interface {
	// Get the current block number
	GetBlockNumber() (uint64, error)
	// GetBlock(number interface{}) (interface{}, error) // number is a digit or 'latest'
	// GetBlockByNumber(number interface{}, fullTx bool) (interface{}, error)
	// // number: "latest"/or number
	// GetBalanceAt(address string, number interface{}) (*big.Int, error)
	// GetBalance(address string) *big.Int
}
