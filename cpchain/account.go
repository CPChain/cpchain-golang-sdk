package cpchain

import (
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

type URL struct {
	Scheme string // Protocol scheme to identify a capable account backend
	Path   string // Path for the backend to identify a unique entity
}

type Account struct {
	Address common.Address `json:"address"`
	URL     URL            `json:"url"`
}

type Wallet interface {
	URL() URL

	Status() (string, error)

	Open(passphrase string) error

	Close() error

	Accounts() []Account

	Contains(account Account) bool

	// Derive(path DerivationPath, pin bool) (Account, error)

	// SelfDerive(base DerivationPath, chain cpchain.ChainStateReader)

	SignHash(account Account, hash []byte) ([]byte, error)

	SignTx(account Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	SignHashWithPassphrase(account Account, passphrase string, hash []byte) ([]byte, error)

	SignTxWithPassphrase(account Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	DecryptWithEcies(account Account, cipherText []byte) ([]byte, error)

	PublicKey(account Account) ([]byte, error)
}
