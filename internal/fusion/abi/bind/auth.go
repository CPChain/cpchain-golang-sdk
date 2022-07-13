package bind

import (
	"crypto/ecdsa"
	"errors"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

func NewKeyedTransactor(key *ecdsa.PrivateKey) *TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}
