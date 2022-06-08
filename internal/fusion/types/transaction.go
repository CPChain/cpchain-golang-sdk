package types

import (
	"math/big"
	"sync/atomic"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
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

func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		// according to eip155, v = chain_id*2 + 35/36
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}



type sigCache struct {
	signer Signer
	from   common.Address
}


// SignTx signs the transaction using the given signer and private key
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}


type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction) (common.Address, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	// Hash returns the hash to be signed.
	Hash(tx *Transaction) common.Hash
	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &Transaction{data: tx.data}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	return cpy, nil
}


type Cep1Signer struct {
	// cf. eip 155
	// v = chain_id*2 + 8 + 27/28
	chainId, chainIdMul *big.Int
}

func NewCep1Signer(chainId *big.Int) Cep1Signer {
	if chainId == nil {
		// chainId = big.NewInt(configs.MainnetChainId)
		chainId = big.NewInt(337) //TODO need change
	}
	return Cep1Signer{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}


// Sender recovers sender address
func (s Cep1Signer) Sender(tx *Transaction) (common.Address, error) {
	if !tx.Protected() {
		// log.Debug("Deprecated signer with unprotected transaction")
		return HomesteadSigner{}.Sender(tx)
	}
	// verify chain id
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}

	V := new(big.Int).Sub(tx.data.V, s.chainIdMul)
	V.Sub(V, big8)

	return recoverPlain(s.Hash(tx), tx.data.R, tx.data.S, V, true)
}

func (s Cep1Signer) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.data.Type,
		tx.data.AccountNonce,
		tx.data.Price,
		tx.data.GasLimit,
		tx.data.Recipient,
		tx.data.Amount,
		tx.data.Payload,
		s.chainId,
		uint(0),
		uint(0),
	})
}

// Signature returns a new transaction with the given signature. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s Cep1Signer) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	R, S, V, err = HomesteadSigner{}.SignatureValues(tx, sig)
	if err != nil {
		return nil, nil, nil, err
	}
	// Sign == 0 only when chainId == 0
	if s.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, s.chainIdMul)
	}
	return R, S, V, nil
}

func (s Cep1Signer) Equal(s2 Signer) bool {
	cep1, ok := s2.(Cep1Signer)
	return ok && cep1.chainId.Cmp(s.chainId) == 0
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		return common.Address{}, ErrInvalidSig
	}
	// v should be 0, 1 by now
	V := byte(Vb.Uint64() - 27)

	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return common.Address{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// HomesteadTransaction implements TransactionInterface using the
// homestead rules.
type HomesteadSigner struct{ FrontierSigner }

func (s HomesteadSigner) Equal(s2 Signer) bool {
	_, ok := s2.(HomesteadSigner)
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (hs HomesteadSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	return hs.FrontierSigner.SignatureValues(tx, sig)
}

func (hs HomesteadSigner) Sender(tx *Transaction) (common.Address, error) {
	return recoverPlain(hs.Hash(tx), tx.data.R, tx.data.S, tx.data.V, true)
}

type FrontierSigner struct{}

func (s FrontierSigner) Equal(s2 Signer) bool {
	_, ok := s2.(FrontierSigner)
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (fs FrontierSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) != 65 {
		panic(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (fs FrontierSigner) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.data.AccountNonce,
		tx.data.Price,
		tx.data.GasLimit,
		tx.data.Recipient,
		tx.data.Amount,
		tx.data.Payload,
	})
}

func (fs FrontierSigner) Sender(tx *Transaction) (common.Address, error) {
	return recoverPlain(fs.Hash(tx), tx.data.R, tx.data.S, tx.data.V, false)
}

// block.go
func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}