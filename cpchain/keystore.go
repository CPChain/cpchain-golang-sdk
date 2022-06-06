package cpchain

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore" //TOOD 是否要换成fusion
)

func NewKeystore(password string) {
	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())
}
