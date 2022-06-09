package types

import (
	"math/big"
	"sync/atomic"
	"errors"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/rlp"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto/sha3"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
)


var (
    big8 = big.NewInt(8)
	ErrInvalidChainId = errors.New("invalid chain id for signer")
	ErrInvalidSig     = errors.New("invalid transaction v, r, s values")
)

type Transaction struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type txdata struct {
	// Type indicates the features assigned to the tx, e.g. private tx.
	Type         uint64          `json:"type" gencodec:"required"`
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

func (tx *Transaction) Protected() bool {
	return isProtectedV(tx.data.V)
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		// cf. eip155 for details
		return v != 27 && v != 28
	}
	// anything not 27 or 28 are considered protected
	return true
}

func (tx *Transaction) ChainId() *big.Int {
	return deriveChainId(tx.data.V)
}


// block.go
func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}