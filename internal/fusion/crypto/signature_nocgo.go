package crypto

import (
	"fmt"
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/btcd/btcec"
)
// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given hash cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign(hash []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	if prv.Curve != btcec.S256() {
		return nil, fmt.Errorf("private key curve is not secp256k1")
	}
	sig, err := btcec.SignCompact(btcec.S256(), (*btcec.PrivateKey)(prv), hash, false)
	if err != nil {
		return nil, err
	}
	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, nil
}

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, sig []byte) ([]byte, error) {
	pub, err := SigToPub(hash, sig)
	if err != nil {
		return nil, err
	}
	bytes := (*btcec.PublicKey)(pub).SerializeUncompressed()
	return bytes, err
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	// Convert to btcec input format with 'recovery id' v at the beginning.
	btcsig := make([]byte, 65)
	btcsig[0] = sig[64] + 27
	copy(btcsig[1:], sig)

	pub, _, err := btcec.RecoverCompact(btcec.S256(), btcsig, hash)
	return (*ecdsa.PublicKey)(pub), err
}


// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return btcec.S256()
}