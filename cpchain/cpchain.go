package cpchain

import (
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

	"github.com/zgljl2012/slog"
)

type cpchain struct {
	network  Network
	provider fusion.Provider
	web3     fusion.Web3
}

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

func (c *cpchain) Backend() (bind.ContractBackend, error) {
	backend, err := cpcclient.Dial(c.network.JsonRpcUrl)
	return backend, err
}

func (c *cpchain) Contract(abi []byte, address string) Contract {
	contractIns, err := contract.NewContractWithProvider(
		[]byte(abi),
		common.HexToAddress(address),
		c.provider,
	)
	if err != nil {
		slog.Fatal(err)
	}
	return &contractInternal{
		contractIns: contractIns,
	}
}

func (c *cpchain) LoadWallet(path string) Wallet {
	account := ReadAccount(path) // 获取账户信息
	// walletbkd := backends.NewClientBackend(c.provider) // 创建一个client
	walletbkd, err := cpcclient.Dial(c.network.JsonRpcUrl)
	if err != nil {
		slog.Fatal(err)
	}
	return &WalletInstance{
		account: *account,
		backend: walletbkd,
		network: c.network,
	}
}

func (c *cpchain) CreateAccount(path string, password string) (*Account, error) {
	key, err := newKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	acct := Account{Address: key.Address, URL: URL{Scheme: KeyStoreScheme, Path: filepath.Join(path, keyFileName(key.Address))}} //TODO path 是否是绝对路径的问题
	if err = StoreKey(key, acct, password); err != nil {
		return nil, err
	}
	return &acct, nil
}

func StoreKey(key *Key, acct Account, password string) error { //TODO 是否应该写入接口内
	keyjson, err := EncryptKey(key, password, 2, 1)
	if err != nil {
		return err
	}
	return writeKeyFile(acct.URL.Path, keyjson)
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
