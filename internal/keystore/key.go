package keystore

import (
	"crypto/ecdsa"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/pborman/uuid"
)

const (
	version = 3
)

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

type keyStore interface {
	// Loads and decrypts the key from disk.
	GetKey(addr common.Address, filename string, auth string) (*Key, error)
	// Writes and encrypts the key.


	// StoreKey(filename string, k *Key, auth string) error


	// Joins filename with the key directory unless it is already absolute.


	// JoinPath(filename string) string
}