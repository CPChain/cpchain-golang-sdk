package cpchain

var Mainnet = Network{
	Name:       "Mainnet",
	JsonRpcUrl: "https://civilian.cpchain.io",
	ChainId:    337,
}

var Testnet = Network{
	Name:       "Testnet",
	JsonRpcUrl: "https://civilian.testnet.cpchain.io",
	ChainId:    41,
}
