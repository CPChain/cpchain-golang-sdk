package types

import (
	"math/big"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/crypto/sha3"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/rlp"
)

var (
	big8 = big.NewInt(8)
)

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
