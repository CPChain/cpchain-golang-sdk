package fusion_test

import (
	"io/ioutil"
	"math/big"
	"testing"

	"cpchain-golang-sdk/internal/fusion"
	"cpchain-golang-sdk/internal/fusion/common"
	"cpchain-golang-sdk/internal/fusion/contract"
)

type CreateProductEvent struct {
	Id        *big.Int       `json:"ID"`
	Name      string         `json:"name"`
	Extend    string         `json:"extend"`
	Price     *big.Int       `json:"price"`
	Creator   common.Address `json:"creator"`
	File_uri  string         `json:"file_uri" rlp:"file_uri"`
	File_hash string         `json:"file_hash"`
	// Raw      types.Log      // Blockchain specific contextual infos
}

func TestFilterEvents(t *testing.T) {
	file, err := ioutil.ReadFile("../../fixtures/product.json")
	if err != nil {
		t.Fatal(err)
	}
	address := common.HexToAddress("0x49F431A6bE97bd26bD416D6E6A0D3FAF3E3d5071")

	URL := "https://civilian.cpchain.io"
	provider, err := fusion.NewHttpProvider(URL)
	if err != nil {
		t.Fatal(err)
	}
	instance, err := contract.NewContractWithProvider(file, address, provider)
	if err != nil {
		t.Fatal(err)
	}
	events, err := instance.FilterLogs("CreateProduct", CreateProductEvent{}, contract.WithFilterLogsFromBlock(6712515))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Count:", len(events))
	for _, e := range events {
		args := e.Data.(*CreateProductEvent)
		t.Log(e.BlockNumber, args.Id, args.Name, args.Price, args.Extend, args.File_hash, args.File_uri)
	}
}
