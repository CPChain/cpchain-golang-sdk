package crypto

import (
	"crypto/elliptic"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto/secp256k1"
)

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
}
