package utils

import (
	"errors"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
)

func GetNetWork(endpoint string) (cpchain.Network, error) {
	if endpoint == cpchain.Mainnet.JsonRpcUrl {
		return cpchain.Mainnet, nil
	} else if endpoint == cpchain.Testnet.JsonRpcUrl {
		return cpchain.Testnet, nil
	} else {
		return cpchain.Network{}, errors.New("endpoint is error")
	}
}

const Abi = "[{\"inputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"constructor\"}]"

const Bin = `0x6080604052348015600f57600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550603f80605d6000396000f3fe6080604052600080fdfea2646970667358221220cc46356d887799b33b3ca82fcf610da45d06ecf8fa0e763740abfbd51f6898ff64736f6c634300080a0033`

func ReadContract(path string) (string, string, error) {
	return Abi, Bin, nil
}
