package cpchain_test

import (
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
)

func TestMock(t *testing.T) {
	e := cpchain.MockEvent("CreateProduct", map[string]interface{}{
		"name": "1",
	}, "0x01", 1)
	t.Log(e.TxHash.Hex())
	t.Log(e.BlockHash.Hex())
	t.Log(e.BlockNumber)
	t.Log(e.Data)
}
