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
	"github.com/CPChain/cpchain-golang-sdk/internal/keystore"
	"github.com/zgljl2012/slog"
)

type WalletInstance struct {
	account Account              // address and url
	backend bind.ContractBackend // client
	key     *keystore.Key        // store key
	network Network              //TODO 只是为了获取chainid
}

// return the address of wallet
func (w *WalletInstance) Addr() common.Address {
	return w.account.Address
}

// return the key of wallet
func (w *WalletInstance) Key() *keystore.Key {
	return w.key
}

// sign transaction with tx and chainid
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

// transfer to target address, if success, return the signtx
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
		slog.Fatal("Estimate gaslimit failed", "err", err)
		return nil, err
	}

	tx := types.NewTransaction(nonce, to, valueInCpc, gasLimit, gasPrice, nil)

	chainID := big.NewInt(0).SetUint64(uint64(w.network.ChainId))

	signedTx, err := w.SignTx(tx, chainID)
	if err != nil {
		slog.Fatal("Sign tx failed", "err", err)
		return nil, err
	}

	err = w.backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		slog.Fatal("Send transaction failed", "err", err)
		return nil, err
	}

	return signedTx, nil
}

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
