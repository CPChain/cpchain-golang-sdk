package cpchain

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/CPChain/cpchain-golang-sdk/internal/cpcclient"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/abi/bind"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/contract"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
	"github.com/CPChain/cpchain-golang-sdk/internal/keystore"
)

type cpchain struct {
	network  Network
	provider fusion.Provider
	web3     fusion.Web3
}

// get network according to the endpoint
func GetNetWork(endpoint string) (Network, error) {
	if endpoint == Mainnet.JsonRpcUrl {
		return Mainnet, nil
	} else if endpoint == Testnet.JsonRpcUrl {
		return Testnet, nil
	} else {
		return Network{}, errors.New("endpoint is error")
	}
}

func NewCPChain(network Network) (CPChain, error) {
	provider, err := fusion.NewHttpProvider(network.JsonRpcUrl)
	if err != nil {
		return nil, err
	}
	web3, err := fusion.NewWeb3(provider)
	if err != nil {
		return nil, err
	}
	return &cpchain{
		provider: provider,
		network:  network,
		web3:     web3,
	}, nil
}

func (c *cpchain) BlockNumber() (uint64, error) {
	block, err := c.web3.GetBlock("latest")
	if err != nil {
		return 0, fmt.Errorf("get block number failed: %v", err)
	}
	return block.(*fusion.FullBlock).Number, nil
}

func (c *cpchain) Block(number int) (*fusion.FullBlock, error) {
	block, err := c.web3.GetBlock(uint64(number))
	if err != nil {
		return nil, fmt.Errorf("get block failed: %v", err)
	}
	return block.(*fusion.FullBlock), nil
}

func (c *cpchain) BalanceOf(address string) *big.Int {
	balance, err := c.web3.GetBalanceAt(address, "latest")
	if err != nil {
		return big.NewInt(0)
	}
	return balance
}

func WeiToCpc(wei *big.Int) *big.Int {
	return wei.Div(wei, big.NewInt(1e18))
}

// create a contract instance
func (c *cpchain) Contract(abi []byte, address string) Contract {
	// contractIns, err := contract.NewContractWithProvider(
	// 	[]byte(abi),
	// 	common.HexToAddress(address),
	// 	c.provider,
	// )
	contractIns, err := contract.NewContractWithUrl(
		[]byte(abi),
		common.HexToAddress(address),
		c.network.JsonRpcUrl,
	)
	if err != nil {
		return nil //TODO 错误处理
	}
	return &contractInternal{
		contractIns: contractIns,
	}
}

// create a wallet instance by import keystorefile and password
func (c *cpchain) LoadWallet(path string, password string) (Wallet, error) {
	account, err := ReadAccount(path) // get account info from file
	if err != nil {
		return nil, fmt.Errorf("load wallet failed: %v", err)
	}
	// walletbkd := backends.NewClientBackend(c.provider) // create a client
	walletbkd, err := cpcclient.Dial(c.network.JsonRpcUrl)
	if err != nil {
		return nil, err
	}
	key, err := GetKey(path, account.Address, password) //get key(include privatekey), need password
	if err != nil {
		return nil, err
	}
	return &WalletInstance{
		account: *account,
		backend: walletbkd,
		network: c.network,
		key:     key,
	}, nil
}

// create a new account on the chain. return the account instance(include address and keystorefile path), and generate a keystorefile in the path where you want to store it
func (c *cpchain) CreateAccount(path string, password string) (*Account, error) {
	pathabs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	key, err := keystore.NewKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	acct := Account{Address: key.Address, URL: URL{Scheme: KeyStoreScheme, Path: filepath.Join(pathabs, keystore.KeyFileName(key.Address))}} // create account instance
	if err = StoreKey(key, acct, password); err != nil {                                                                                     // create keystorefile
		return nil, err
	}
	return &acct, nil
}

// Deploy contract,get contract abi and bin from file
func (c *cpchain) DeployContractByFile(path string, w Wallet) (common.Address, *types.Transaction, error) {
	abi, bin, err := ReadContract(path)
	if err != nil {
		return common.Address{}, nil, err
	}
	return c.DeployContract(abi, bin, w)
}

// Deploy contract, import a wallet instance to send this transaction
func (c *cpchain) DeployContract(abi string, bin string, w Wallet) (common.Address, *types.Transaction, error) {
	key := w.Key()                                       // get key
	backend, err := cpcclient.Dial(c.network.JsonRpcUrl) // create backend
	if err != nil {
		return common.Address{}, nil, err
	}

	nonce, err := backend.PendingNonceAt(context.Background(), w.Addr()) //get nonce
	if err != nil {
		return common.Address{}, nil, err
	}

	auth := contract.NewTransactor(key.PrivateKey, new(big.Int).SetUint64(nonce))

	// // address, tx, contract, err := contract.DeployContract(abi, auth, common.FromHex(bin), w.backend, w.network.ChainId)
	address, tx, _, err := contract.DeployContract(abi, auth, common.FromHex(bin), backend, c.network.ChainId)
	if err != nil {
		return common.Address{}, nil, nil
	}
	return address, tx, nil
}

func (c *cpchain) ReceiptByTx(signedTx *types.Transaction) (*types.Receipt, error) {
	backend, err := cpcclient.Dial(c.network.JsonRpcUrl)
	if err != nil {
		return &types.Receipt{}, err
	}
	receipt, err := bind.WaitMined(context.Background(), backend, signedTx)
	if err != nil {
		return &types.Receipt{}, err
	}
	return receipt, nil
}

// Write an account and its corresponding key to the keystore file
func StoreKey(key *keystore.Key, acct Account, password string) error { //TODO 是否应该写入接口内
	keyjson, err := keystore.EncryptKey(key, password, 2, 1)
	if err != nil {
		return err
	}
	return keystore.WriteKeyFile(acct.URL.Path, keyjson)
}

type contractInternal struct {
	contractIns contract.Contract
}

func (c *contractInternal) Events(eventName string, event interface{}, options ...WithEventsOptionsOption) ([]*contract.Event, error) {
	var opts = EventsOptions{}
	for _, op := range options {
		op(&opts)
	}
	events, err := c.contractIns.FilterLogs(
		eventName,
		event,
		contract.WithFilterLogsFromBlock(opts.FromBlock),
		contract.WithFilterLogsEndBlock(opts.ToBlock),
	)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c *contractInternal) Call(w Wallet, chainId uint, method string, params ...interface{}) (*types.Transaction, error) {
	key := w.Key()
	backend, err := cpcclient.Dial(Testnet.JsonRpcUrl) //TODO
	// Key, err := w.GetKey(w)
	if err != nil {
		return nil, err
	}

	nonce, err := backend.PendingNonceAt(context.Background(), w.Addr())
	if err != nil {
		return nil, err
	}

	auth := contract.NewTransactor(key.PrivateKey, new(big.Int).SetUint64(nonce))

	tx, err := c.contractIns.Transact(auth, chainId, method, params...)
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (c *contractInternal) View(result interface{}, method string, params ...interface{}) error {
	callOpts := NewCallOpt()
	err := c.contractIns.Call(callOpts, result, method, params...)
	return err
}

func NewCallOpt() *bind.CallOpts {
	return &bind.CallOpts{
		Pending: false,
		From:    common.Address{},
		Context: context.Background(),
	}
}
