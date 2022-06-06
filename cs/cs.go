package main

import (
	"fmt"
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
)

func main() {
	clientOnMainnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		fmt.Println("------")
	}
	balanceOnMainnet := clientOnMainnet.BalanceOf("0xc0d69e165Bbc8a4ce33107bB9768f55587162B16")
	if balanceOnMainnet.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("------")
	}
	fmt.Println(cpchain.WeiToCpc(balanceOnMainnet))
}
