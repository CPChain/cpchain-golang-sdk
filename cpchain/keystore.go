package cpchain

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto"
)

func NewKeyStore(keydir string, scryptN, scryptP int) *KeyStore {
	keydir, _ = filepath.Abs(keydir)
	ks := &KeyStore{storage: &keyStorePassphrase{keydir, scryptN, scryptP}}
	ks.account = ReadAccount(keydir)
	return ks
}

func DecryptKeystore(keyStoreFilePath string, password string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address, *KeyStore, Account, error) {
	// Open keystore file.
	file, err := os.Open(keyStoreFilePath)
	if err != nil {
		return nil, nil, [20]byte{}, nil, Account{}, err
	}
	keyPath, err := filepath.Abs(filepath.Dir(file.Name()))
	if err != nil {
		return nil, nil, [20]byte{}, nil, Account{}, err
	}
	// Create keystore and get account.
	kst := NewKeyStore(keyPath, 2, 1)
	acct := kst.Accounts()
	fmt.Println("-----")
	fmt.Println(acct)
	fmt.Println("-----")
	key, err := kst.storage.GetKey(acct.Address, acct.URL.Path, password)
	fmt.Println("-----")
	if err != nil {
		return nil, nil, [20]byte{}, nil, Account{}, err
	}
	// Get private and public keys.
	privateKey := key.PrivateKey
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, [20]byte{}, nil, Account{}, errors.New("error casting public key to ECDSA")
	}

	// Get contractAddress.
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, publicKeyECDSA, fromAddress, kst, *acct, nil

}
