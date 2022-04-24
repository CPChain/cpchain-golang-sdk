package cpchain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/common"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/types"
)

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func MockEvent(name string, data interface{}, address string,
	blockNumber int) Event {
	return Event{
		Name: name,
		Data: data,
		Log: types.Log{
			Address:     common.HexToAddress(address),
			BlockNumber: uint64(blockNumber),
			TxHash:      common.HexToHash("0x" + randomString(62)),
			BlockHash:   common.HexToHash("0x" + randomString(62)),
		},
	}
}
