package cpchain

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	// "math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto"
)

func NewKeyStore(keydir string, scryptN, scryptP int) *KeyStore {
	keydir, _ = filepath.Abs(keydir)
	ks := &KeyStore{storage: &keyStorePassphrase{keydir, scryptN, scryptP}}
	ks.account = ReadAccount(keydir)
	return ks
}


// var (
// 	endPoint = "http://localhost:8501"
// )


func resolveDomain(hostname string) (string, error) {
	ipAddress := net.ParseIP(hostname)
	// log.Debug("parse ip", "hostname", hostname, "ipAddress", ipAddress)
	if ipAddress != nil {
		return ipAddress.String(), nil
	}
	addr, err := net.LookupHost(hostname)
	if err != nil {
		// log.Error("lookup host error", "hostname", hostname, "err", err)
		return "", err
	}
	if len(addr) > 0 {
		return addr[0], nil
	}
	return "", fmt.Errorf("invalid host: %v", err)
}

func ResolveUrl(url string) (string, error) {
	host, port, err := net.SplitHostPort(url[7:])
	ip, err := resolveDomain(host)
	if err != nil {
		// log.Fatal("unknown endpoint", "endpoint", url, "err", err)
		return "", err
	}
	return "http://" + ip + ":" + port, err
}

func Connect(endPoint string, keyStoreFilePath string, password string) (*Client, error, *ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address, *KeyStore, Account) {
	// ep, err := ResolveUrl(endPoint)
	// fmt.Println(ep)
	// if err != nil {

	// }
	// Create client.
	client, err := Dial(endPoint)
	if err != nil {
		// log.Fatal(err.Error())
	}

	privateKey, publicKeyECDSA, fromAddress, kst, account, err := DecryptKeystore(keyStoreFilePath, password)
	if err != nil {
		// log.Fatal(err.Error())
	}
	return client, err, privateKey, publicKeyECDSA, fromAddress, kst, account
}




func DecryptKeystore(keyStoreFilePath string, password string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address, *KeyStore, Account, error) {
	// Open keystore file.
	_, err := os.Open(keyStoreFilePath)
	if err != nil {
		return nil, nil, [20]byte{}, nil, Account{}, err
	}
	// Create keystore and get account.
	kst := NewKeyStore(keyStoreFilePath, 2, 1)
	acct := kst.Accounts()
	key, err := kst.storage.GetKey(acct.Address, acct.URL.Path, password)
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
