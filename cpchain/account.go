package cpchain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

type Account struct {
	Address common.Address `json:"address"` // cpchain account address derived from the key
	URL     URL            `json:"url"`     // Optional resource locator within a backend
}

type URL struct {
	Path   string
	Scheme string
}

func (u URL) Cmp(url URL) int {
	if u.Scheme == url.Scheme {
		return strings.Compare(u.Path, url.Path)
	}
	return strings.Compare(u.Scheme, url.Scheme)
}

func (a *Account) Addr() common.Address {
	return a.Address
}

func (a *Account) GetKey(password string) (*Key, error) {
	// Load the key from the keystore and decrypt its contents
	keyjson, err := ioutil.ReadFile(a.URL.Path)
	if err != nil {
		return nil, err
	}
	key, err := DecryptKey(keyjson, password)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != a.Address {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, a.Address)
	}
	return key, nil
}

func (a *Account) SignTxWithPassword(password string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	key, err := a.GetKey(password)
	if err != nil {
		return nil, err
	}
	privateKey := key.PrivateKey
	signTx, err := types.SignTx(tx, types.NewCep1Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}
	return signTx, nil
}

var (
	buf  = new(bufio.Reader)
	keys struct {
		Address string `json:"address"`
	}
)

const KeyStoreScheme = "keystore"

func ReadAccount(path string) (*Account, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open keystore file, path: %s, err: %v", path, err)
	}
	defer fd.Close()
	buf.Reset(fd)
	// Parse the address.
	keys.Address = ""
	err = json.NewDecoder(buf).Decode(&keys)
	addr := common.HexToAddress(keys.Address)
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to decode keystore key, path: %s, err: %v", path, err)
	case (addr == common.Address{}):
		return nil, fmt.Errorf("failed to decode keystore key, path: %s, err: missing or zero address", path)
	default:
		return &Account{Address: addr, URL: URL{Scheme: KeyStoreScheme, Path: path}}, nil
	}
}

func GetKey(path string, account common.Address, password string) (*Key, error) {
	// Load the key from the keystore and decrypt its contents
	keyjson, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := DecryptKey(keyjson, password)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != account {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, account)
	}
	return key, nil
}
