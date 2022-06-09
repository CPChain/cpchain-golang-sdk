package cpchain

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"
	"errors"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/accounts"
	"github.com/CPChain/cpchain-golang-sdk/internal/keystore"
)


// func Connect(password string) (*cpclient.Client, error, *ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address, *keystore.KeyStore, accounts.Account, *big.Int) {
// 	ep, err := configs.ResolveUrl(endPoint)
// 	if err != nil {
// 		// log.Fatal("unknown endpoint", "endpoint", endPoint, "err", err)
// 	}
// 	// Create client.
// 	client, err := cpclient.Dial(ep)
// 	if err != nil {
// 		// log.Fatal(err.Error())
// 	}

// 	chainConfig, err := client.ChainConfig()
// 	if err != nil {
// 		// log.Fatal(err.Error())
// 	}
// 	chainId, runMode := chainConfig.ChainID.Uint64(), configs.Mainnet
// 	switch chainId {
// 	case configs.DevChainId:
// 		runMode = configs.Dev
// 	case configs.MainnetChainId:
// 		runMode = configs.Mainnet
// 	case configs.TestMainnetChainId:
// 		runMode = configs.TestMainnet
// 	case configs.TestnetChainId:
// 		runMode = configs.Testnet
// 	default:
// 		log.Fatal("unknown chain id")
// 	}
// 	configs.SetRunMode(runMode)

// 	privateKey, publicKeyECDSA, fromAddress, kst, account, err := DecryptKeystore(password)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	return client, err, privateKey, publicKeyECDSA, fromAddress, kst, account, big.NewInt(0).SetUint64(chainId)
// }


func DecryptKeystore(keyStoreFilePath string, password string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address, *keystore.KeyStore, accounts.Account, error) {
	// Open keystore file.
	file, err := os.Open(keyStoreFilePath)
	if err != nil {
		return nil, nil, [20]byte{}, nil, accounts.Account{}, err
	}
	keyPath, err := filepath.Abs(filepath.Dir(file.Name()))
	if err != nil {
		return nil, nil, [20]byte{}, nil, accounts.Account{}, err
	}
	// Create keystore and get account.
	kst := keystore.NewKeyStore(keyPath, 2, 1)
	account := kst.Accounts()[0]
	account, key, err := kst.GetDecryptedKey(account, password)
	if err != nil {
		return nil, nil, [20]byte{}, nil, accounts.Account{}, err
	}
	// Get private and public keys.
	privateKey := key.PrivateKey
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, [20]byte{}, nil, accounts.Account{}, errors.New("error casting public key to ECDSA")
	}

	// Get contractAddress.
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, publicKeyECDSA, fromAddress, kst, account, nil

}