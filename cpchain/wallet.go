package cpchain

import (
	"bufio"
	"context"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"

	"github.com/CPChain/cpchain-golang-sdk/internal/cpcclient"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/contract"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
	"github.com/zgljl2012/slog"
)

type WalletInstance struct {
	account Account
	backend bind.ContractBackend
	key     *Key
	network Network //TODO 只是为了获取chainid
}

func (w *WalletInstance) Addr() common.Address {
	return w.account.Address
}

func (w *WalletInstance) Key() *Key {
	return w.key
}

// func (w *WalletInstance) GetKey(password string) (*Key, error) {
// 	// Load the key from the keystore and decrypt its contents
// 	keyjson, err := ioutil.ReadFile(w.account.URL.Path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	key, err := DecryptKey(keyjson, password)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Make sure we're really operating on the requested key (no swap attacks)
// 	if key.Address != w.account.Address {
// 		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, w.account.Address)
// 	}
// 	return key, nil
// }

// func (w *WalletInstance) SignTxWithPassword(password string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
// 	key, err := w.GetKey(password)
// 	if err != nil {
// 		return nil, err
// 	}
// 	privateKey := key.PrivateKey
// 	signTx, err := types.SignTx(tx, types.NewCep1Signer(chainID), privateKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return signTx, nil
// }

func (w *WalletInstance) SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	privateKey := w.key.PrivateKey
	signTx, err := types.SignTx(tx, types.NewCep1Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}
	return signTx, nil
}

const (
	Cpc = 1e18
)

func (w *WalletInstance) Transfer(targetAddr string, value int64) (*types.Transaction, error) {
	fromAddr := w.Addr()

	nonce, err := w.backend.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		slog.Fatal("Get nonce failed: %v", err)
		return nil, err
	}

	gasPrice, err := w.backend.SuggestGasPrice(context.Background())
	if err != nil {
		slog.Fatal("Get gasprice failed: %v", err)
		return nil, err
	}

	to := common.HexToAddress(targetAddr)

	valueInCpc := new(big.Int).Mul(big.NewInt(value), big.NewInt(Cpc))

	msg := cpcclient.CallMsg{From: fromAddr, To: &to, Value: valueInCpc, Data: nil} //TODO

	gasLimit, err := w.backend.EstimateGas(context.Background(), msg)
	if err != nil {
		slog.Fatal("Get gaslimit failed: %v", err)
		return nil, err
	}

	tx := types.NewTransaction(nonce, to, valueInCpc, gasLimit, gasPrice, nil)

	chainID := big.NewInt(0).SetUint64(uint64(w.network.ChainId))

	signedTx, err := w.SignTx(tx, chainID)
	if err != nil {
		slog.Fatal("Sign tx failed: %v", err)
		return nil, err
	}

	err = w.backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		slog.Fatal("Send transaction failed: %v", err)
		return nil, err
	}

	return signedTx, nil
}

// func (w *WalletInstance) DeployContractByFile(path string, password string) (common.Address, *types.Transaction, error) {
// 	abi, bin, err := ReadContract(path)
// 	if err != nil {
// 		slog.Fatal(err)
// 	}
// 	return w.DeployContract(abi, bin, password)
// }

// func (w *WalletInstance) DeployContract(abi string, bin string, password string) (common.Address, *types.Transaction, error) {
// 	Key, err := w.GetKey(password)
// 	if err != nil {
// 		slog.Fatal(err)
// 		return common.Address{}, nil, nil
// 	}

// 	nonce, err := w.backend.PendingNonceAt(context.Background(), w.Addr())
// 	if err != nil {
// 		slog.Fatal(err)
// 		return common.Address{}, nil, nil
// 	}

// 	auth := contract.NewTransactor(Key.PrivateKey, new(big.Int).SetUint64(nonce))

// 	// address, tx, contract, err := contract.DeployContract(abi, auth, common.FromHex(bin), w.backend, w.network.ChainId)
// 	address, tx, _, err := contract.DeployContract(abi, auth, common.FromHex(bin), w.backend, w.network.ChainId)
// 	if err != nil {
// 		return common.Address{}, nil, nil
// 	}
// 	return address, tx, nil
// }

func (w *WalletInstance) NewTransactor(password string) *bind.TransactOpts {
	Key := w.key

	nonce, err := w.backend.PendingNonceAt(context.Background(), w.Addr())
	if err != nil {
	}

	auth := contract.NewTransactor(Key.PrivateKey, new(big.Int).SetUint64(nonce))

	return auth
}

var (
	bufs         = new(bufio.Reader)
	contractjson struct {
		Abi interface{} `json:"abi"`
		Bin string      `json:"bytecode"`
	}
)

func ReadContract(path string) (string, string, error) {
	fpath, err := filepath.Abs(path)
	if err != nil {
		slog.Fatal(err)
	}
	contractFile, err := os.Open(fpath)
	if err != nil {
		slog.Fatal(err)
	}
	defer contractFile.Close()
	bufs.Reset(contractFile)
	// Parse the address.
	err = json.NewDecoder(bufs).Decode(&contractjson)
	if err != nil {
		return "", "", nil
	}
	abijson, err := json.Marshal(contractjson.Abi)
	if err != nil {
		return "", contractjson.Bin, nil
	}
	abistring := string(abijson)

	return abistring, contractjson.Bin, err
}
