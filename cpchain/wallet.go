package cpchain

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"

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
	network Network //TODO 只是为了获取chainid
}

func (w *WalletInstance) Addr() common.Address {
	return w.account.Address
}

func (w *WalletInstance) GetKey(password string) (*Key, error) {
	// Load the key from the keystore and decrypt its contents
	keyjson, err := ioutil.ReadFile(w.account.URL.Path)
	if err != nil {
		return nil, err
	}
	key, err := DecryptKey(keyjson, password)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != w.account.Address {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, w.account.Address)
	}
	return key, nil
}

func (w *WalletInstance) SignTxWithPassword(password string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	key, err := w.GetKey(password)
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

const (
	Cpc = 1e18
)

func (w *WalletInstance) Transfer(password string, targetAddr string, value int64) error {
	fromAddr := w.Addr()

	nonce, err := w.backend.PendingNonceAt(context.Background(), fromAddr)

	gasPrice, err := w.backend.SuggestGasPrice(context.Background())

	to := common.HexToAddress(targetAddr)

	valueInCpc := new(big.Int).Mul(big.NewInt(value), big.NewInt(Cpc))

	msg := cpcclient.CallMsg{From: fromAddr, To: &to, Value: valueInCpc, Data: nil} //TODO

	gasLimit, err := w.backend.EstimateGas(context.Background(), msg)

	tx := types.NewTransaction(nonce, to, valueInCpc, gasLimit, gasPrice, nil)

	chainID := big.NewInt(0).SetUint64(uint64(w.network.ChainId))

	signedTx, err := w.SignTxWithPassword(password, tx, chainID)

	err = w.backend.SendTransaction(context.Background(), signedTx)

	if err != nil {
		return err
	}
	return nil
}

func (w *WalletInstance) DeployContractByFile(path string, password string) error {
	abi, bin, err := ReadContract(path)
	if err != nil {
		slog.Fatal(err)
	}
	return w.DeployContract(abi, bin, password)
}

func (w *WalletInstance) DeployContract(abi string, bin string, password string) error {
	Key, err := w.GetKey(password)
	if err != nil {
		slog.Fatal(err)
	}

	nonce, err := w.backend.PendingNonceAt(context.Background(), w.Addr())
	if err != nil {
		slog.Fatal(err)
	}

	auth := contract.NewTransactor(Key.PrivateKey, new(big.Int).SetUint64(nonce))

	// address, tx, contract, err := contract.DeployContract(abi, auth, common.FromHex(bin), w.backend, w.network.ChainId)
	_, _, _, err = contract.DeployContract(abi, auth, common.FromHex(bin), w.backend, w.network.ChainId)
	if err != nil {
		return nil
	}
	return nil
}

const Abi = "[{\"inputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"constructor\"}]"

const Bin = `0x6080604052348015600f57600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550603f80605d6000396000f3fe6080604052600080fdfea2646970667358221220cc46356d887799b33b3ca82fcf610da45d06ecf8fa0e763740abfbd51f6898ff64736f6c634300080a0033`

func ReadContract(path string) (string, string, error) {
	return Abi, Bin, nil
}
