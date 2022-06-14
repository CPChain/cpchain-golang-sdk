package cpchain

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"context"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/contract"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto/ecies"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/rpc"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common/hexutil"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/rlp"
	"github.com/pborman/uuid"
)

type Event = contract.Event

type Network struct {
	Name       string
	JsonRpcUrl string
	ChainId    uint
}

type EventsOptions struct {
	FromBlock uint64
	ToBlock   uint64
}

type WithEventsOptionsOption func(*EventsOptions)

func WithEventsOptionsFromBlock(block uint64) WithEventsOptionsOption {
	return func(flo *EventsOptions) {
		flo.FromBlock = block
	}
}

func WithEventsOptionsToBlock(block uint64) WithEventsOptionsOption {
	return func(flo *EventsOptions) {
		flo.ToBlock = block
	}
}

type Contract interface {
	// TODO 如果事件非常多，如10000条事件，是否会分批获取？
	Events(eventName string, event interface{}, options ...WithEventsOptionsOption) ([]*contract.Event, error)
}

type CPChain interface {
	// Get the current block number
	BlockNumber() (uint64, error)
	Block(number int) (*fusion.FullBlock, error)
	// get balance
	BalanceOf(address string) *big.Int
	Contract(abi []byte, address string) Contract
}

// TODO simulate chain

type Account struct {
	Address common.Address `json:"address"` // cpchain account address derived from the key
	URL     URL            `json:"url"`     // Optional resource locator within a backend
}

type URL struct {
	Scheme string // Protocol scheme to identify a capable account backend
	Path   string // Path for the backend to identify a unique entity
}

func (u URL) Cmp(url URL) int {
	if u.Scheme == url.Scheme {
		return strings.Compare(u.Path, url.Path)
	}
	return strings.Compare(u.Scheme, url.Scheme)
}

type Wallet interface {
	URL() URL

	Account() []Account

	SignTxWithPassphrase(account Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}

type KeyStore struct {
	storage keyStore // Storage backend, might be cleartext or encrypted
	// cache    *accountCache                // In-memory account cache over the filesystem storage
	account *Account
}

var (
	buf  = new(bufio.Reader)
	keys struct {
		Address string `json:"address"`
	}
)

const KeyStoreScheme = "keystore"

func ReadAccount(path string) *Account {
	fd, err := os.Open(path)
	if err != nil {
		// log.Debug("Failed to open keystore file", "path", path, "err", err)
		return nil
	}
	defer fd.Close()
	buf.Reset(fd)
	// Parse the address.
	keys.Address = ""
	err = json.NewDecoder(buf).Decode(&keys)
	addr := common.HexToAddress(keys.Address)
	switch {
	case err != nil:
		// log.Debug("Failed to decode keystore key", "path", path, "err", err)
	case (addr == common.Address{}):
		// log.Debug("Failed to decode keystore key", "path", path, "err", "missing or zero address")
	default:
		return &Account{Address: addr, URL: URL{Scheme: KeyStoreScheme, Path: path}}
	}
	return nil
}

func (ks *KeyStore) Accounts() *Account {
	return ks.account
}

type keyStore interface {
	GetKey(addr common.Address, filename string, auth string) (*Key, error)
}

type Key struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey *ecdsa.PrivateKey

	// we cache ecies.PrivateKey used to decrypt data that encrypted by ecies.PublicKey from remote node
	EciesPrivateKey *ecies.PrivateKey
}

type keyStorePassphrase struct {
	keysDirPath string
	scryptN     int
	scryptP     int
}

func (ks keyStorePassphrase) GetKey(addr common.Address, filename, auth string) (*Key, error) {
	// Load the key from the keystore and decrypt its contents
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	key, err := DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, addr)
	}
	return key, nil
}


// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	c *rpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	return DialContext(context.Background(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c *rpc.Client) *Client {
	return &Client{c}
}

func (c *Client) Close() {
	c.c.Close()
}


// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (c *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := c.c.CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	return uint64(result), err
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var hex hexutil.Big
	if err := c.c.CallContext(ctx, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (c *Client) EstimateGas(ctx context.Context, msg CallMsg) (uint64, error) {
	var hex hexutil.Uint64
	err := c.c.CallContext(ctx, &hex, "eth_estimateGas", toCallArg(msg))
	if err != nil {
		return 0, err
	}
	return uint64(hex), nil
}

// SendTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return fmt.Errorf("encode to bytes error: %v", err)
	}
	fmt.Println("--->>>", common.ToHex(data))
	return c.c.CallContext(ctx, nil, "eth_sendRawTransaction", common.ToHex(data))
}

func toCallArg(msg CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

type CallMsg struct {
	From     common.Address  // the sender of the 'transaction'
	To       *common.Address // the destination contract (nil for contract creation)
	Gas      uint64          // if 0, the call executes with near-infinite gas
	GasPrice *big.Int        // wei <-> gas exchange ratio
	Value    *big.Int        // amount of wei sent along with the call
	Data     []byte          // input data, usually an ABI-encoded contract method invocation
}