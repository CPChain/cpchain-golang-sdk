package cpchain_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
)

func TestWalletTransfer(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet, _ := clientOnTestnet.LoadWallet(keystorePath, password)

	tx, err := wallet.Transfer(targetAddr, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Tx hash: %v", tx.Hash().Hex())
}

func TestIsInt(t *testing.T) {
	a := IsInt("1235565.456")
	if a {
		fmt.Println("ok")
	}
}

func IsInt(p string) bool {
	if len(p) >= 2 && string(p[0]) == "0" {
		return false
	} else {
		patterns := regexp.MustCompile(`\d{1,}`)
		result := patterns.FindString(p)
		if result == p {
			return true
		}
		return false
	}
}

func TestInput(t *testing.T) {
	Abi, _, _ := cpchain.ReadContract("../fixtures/contract/Hello.json")
	inter, err := cpchain.ConvertParmas(Abi, "helloToSomeOne", []string{"41"})
	fmt.Println(inter[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(inter)
}
