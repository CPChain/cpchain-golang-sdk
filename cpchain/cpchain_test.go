package cpchain_test

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/CPChain/cpchain-golang-sdk/cpchain"
	"github.com/CPChain/cpchain-golang-sdk/internal/cpcclient"
	"github.com/CPChain/cpchain-golang-sdk/internal/fusion/contract"
)

func TestGetBlockNumber(t *testing.T) {
	clientOnMainnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	blockNumberOnMainnet, err := clientOnMainnet.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if blockNumberOnMainnet == 0 {
		t.Fatal("BlockNumber is 0")
	}
	block1, err := clientOnMainnet.Block(1)
	if err != nil {
		t.Fatal(err)
	}
	if block1.Number != 1 {
		t.Fatal("BlockNumber is error")
	}
}

func TestGetBalance(t *testing.T) {
	clientOnMainnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	balanceOnMainnet := clientOnMainnet.BalanceOf("0x0a1ea332c4d3d457f17e0ada059f7275b3e2ea1e")
	if balanceOnMainnet.Cmp(big.NewInt(0)) == 0 {
		t.Fatal("Balance is 0")
	}
	t.Log(cpchain.WeiToCpc(balanceOnMainnet))
}

// 测试合约的事件
type CreateProductEvent struct {
	Id        cpchain.UInt256 `json:"ID"`
	Name      cpchain.String  `json:"name"`
	Extend    cpchain.String  `json:"extend"`
	Price     cpchain.UInt256 `json:"price"`
	Creator   cpchain.Address `json:"creator"`
	File_uri  cpchain.String  `json:"file_uri" rlp:"file_uri"`
	File_hash cpchain.String  `json:"file_hash"`
}

func TestEvents(t *testing.T) {
	client, err := cpchain.NewCPChain(cpchain.Mainnet)
	if err != nil {
		t.Fatal(err)
	}
	file, err := ioutil.ReadFile("../fixtures/product.json")
	if err != nil {
		t.Fatal(err)
	}
	address := "0x49F431A6bE97bd26bD416D6E6A0D3FAF3E3d5071"
	events, err := client.Contract(file, address).Events("CreateProduct",
		CreateProductEvent{},
		cpchain.WithEventsOptionsFromBlock(6712515))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Count:", len(events))
	for _, e := range events {
		args := e.Data.(*CreateProductEvent)
		t.Log(e.BlockNumber, args.Id, args.Name, args.Price, args.Extend, args.File_hash, args.File_uri, args.Creator.Hex())
		// check event name
		if e.Name != "CreateProduct" {
			t.Fatal("event name is error")
		}
	}
}

func TestCreateWallet(t *testing.T) {
	password := "123456"
	client, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	path, err := ioutil.TempDir("e:/chengtcode/cpchain-golang-sdk/fixtures", "keystore")
	a, err := client.CreateWallet(path, password)
	if err != nil {
		t.Fatal(err)
	}
	w := client.LoadWallet(a.URL.Path)
	key, err := w.GetKey(password)
	if key.Address != a.Address {
		t.Fatal("account error")
	}
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll(path)
}

// RnodeABI is the input ABI used to generate the binding from.
const Abi = "[{\"constant\":true,\"inputs\":[],\"name\":\"getRnodeNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setPeriod\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"quitRnode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isContract\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"enabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"enableContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"refundAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"Participants\",\"outputs\":[{\"name\":\"lockedDeposit\",\"type\":\"uint256\"},{\"name\":\"lockedTime\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"setSupportedVersion\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"disableContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"setRnodeThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isRnode\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"joinRnode\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rnodeThreshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"supportedVersion\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRnodes\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"period\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"investor\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"who\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"lockedDeposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"lockedTime\",\"type\":\"uint256\"}],\"name\":\"NewRnode\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"who\",\"type\":\"address\"}],\"name\":\"RnodeQuit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"who\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ownerRefund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"numOfInvestor\",\"type\":\"uint256\"}],\"name\":\"ownerRefundAll\",\"type\":\"event\"}]"

// RnodeBin is the compiled bytecode used for deploying new contracts.
const Bin = `0x60806040526117706001908155692a5a058fc295ed00000060025560038190556007805460ff1916909117905534801561003857600080fd5b5060008054600160a060020a03191633179055610be88061005a6000396000f3006080604052600436106100fb5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630b443f4281146101005780630f3a9f6514610127578063113c8498146101415780631627905514610156578063238dafe01461018b578063367edd32146101a057806338e771ab146101b5578063595aa13d146101ca5780635f86d4ca14610204578063894ba8331461021c578063975dd4b214610231578063a8f0769714610249578063aae80f781461026a578063b7b3e9da14610275578063d5601e9f1461028a578063e508bb851461029f578063ef78d4fd14610304578063fa89401a14610319575b600080fd5b34801561010c57600080fd5b5061011561033a565b60408051918252519081900360200190f35b34801561013357600080fd5b5061013f600435610341565b005b34801561014d57600080fd5b5061013f61036d565b34801561016257600080fd5b50610177600160a060020a036004351661043e565b604080519115158252519081900360200190f35b34801561019757600080fd5b50610177610446565b3480156101ac57600080fd5b5061013f61044f565b3480156101c157600080fd5b5061013f610475565b3480156101d657600080fd5b506101eb600160a060020a03600435166105c6565b6040805192835260208301919091528051918290030190f35b34801561021057600080fd5b5061013f6004356105df565b34801561022857600080fd5b5061013f6105fb565b34801561023d57600080fd5b5061013f60043561061e565b34801561025557600080fd5b50610177600160a060020a0360043516610651565b61013f60043561066a565b34801561028157600080fd5b506101156107f4565b34801561029657600080fd5b506101156107fa565b3480156102ab57600080fd5b506102b4610800565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156102f05781810151838201526020016102d8565b505050509050019250505060405180910390f35b34801561031057600080fd5b50610115610811565b34801561032557600080fd5b5061013f600160a060020a0360043516610817565b6006545b90565b600054600160a060020a0316331461035857600080fd5b6201518081111561036857600080fd5b600155565b61037e60043363ffffffff61090516565b151561038957600080fd5b600180543360009081526008602052604090209091015442910111156103ae57600080fd5b3360008181526008602052604080822054905181156108fc0292818181858888f193505050501580156103e5573d6000803e3d6000fd5b50336000818152600860205260408120556104089060049063ffffffff61092416565b506040805133815290517f602a2a9c94f70293aa2be9077f0b2dc89d388bc293fdbcd968274f43494c380d9181900360200190a1565b6000903b1190565b60075460ff1681565b600054600160a060020a0316331461046657600080fd5b6007805460ff19166001179055565b60008054819081908190600160a060020a0316331461049357600080fd5b6006549350600092505b83831015610583576006805460009081106104b457fe5b6000918252602080832090910154600160a060020a0316808352600890915260408083205490519194509250839183156108fc02918491818181858888f19350505050158015610508573d6000803e3d6000fd5b50600160a060020a03821660009081526008602052604081205561053360048363ffffffff61092416565b5060408051600160a060020a03841681526020810183905281517f3914ba80eb00486e7a58b91fb4795283df0c5b507eea9cf7c77cce26cc70d25c929181900390910190a160019092019161049d565b6006541561058d57fe5b6040805185815290517fb65ebb6b17695b3a5612c7a0f6f60e649c02ba24b36b546b8d037e98215fdb8d9181900360200190a150505050565b6008602052600090815260409020805460019091015482565b600054600160a060020a031633146105f657600080fd5b600355565b600054600160a060020a0316331461061257600080fd5b6007805460ff19169055565b600054600160a060020a0316331461063557600080fd5b692a5a058fc295ed00000081101561064c57600080fd5b600255565b600061066460048363ffffffff61090516565b92915050565b60075460ff16151561067b57600080fd5b60035481101561068a57600080fd5b6106933361043e565b1561072557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f706c65617365206e6f742075736520636f6e74726163742063616c6c2074686960448201527f732066756e6374696f6e00000000000000000000000000000000000000000000606482015290519081900360840190fd5b61073660043363ffffffff61090516565b1561074057600080fd5b60025434101561074f57600080fd5b3360009081526008602052604090205461076f903463ffffffff610a6a16565b3360008181526008602052604090209182554260019092019190915561079d9060049063ffffffff610a8016565b5033600081815260086020908152604091829020805460019091015483519485529184015282820152517f586bfaa7a657ad9313326c9269639546950d589bd479b3d6928be469d6dc29039181900360600190a150565b60025481565b60035481565b606061080c6004610b0f565b905090565b60015481565b60008054600160a060020a0316331461082f57600080fd5b61084060048363ffffffff61090516565b151561084b57600080fd5b50600160a060020a03811660008181526008602052604080822054905190929183156108fc02918491818181858888f19350505050158015610891573d6000803e3d6000fd5b50600160a060020a0382166000908152600860205260408120556108bc60048363ffffffff61092416565b5060408051600160a060020a03841681526020810183905281517f3914ba80eb00486e7a58b91fb4795283df0c5b507eea9cf7c77cce26cc70d25c929181900390910190a15050565b600160a060020a03166000908152602091909152604090205460ff1690565b600160a060020a03811660009081526020839052604081205481908190819060ff1615156109555760009350610a61565b600160a060020a038516600090815260208781526040808320805460ff1916905560028901805460018b019093529220549094509250600019840184811061099957fe5b600091825260209091200154600287018054600160a060020a0390921692508291849081106109c457fe5b6000918252602080832091909101805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03948516179055918316815260018801909152604090208290556002860180546000198501908110610a2057fe5b6000918252602090912001805473ffffffffffffffffffffffffffffffffffffffff1916905560028601805490610a5b906000198301610b75565b50600193505b50505092915050565b600082820183811015610a7957fe5b9392505050565b600160a060020a03811660009081526020839052604081205460ff1615610aa957506000610664565b50600160a060020a0316600081815260208381526040808320805460ff19166001908117909155600286018054968201845291842086905585810182559083529120909201805473ffffffffffffffffffffffffffffffffffffffff1916909117905590565b606081600201805480602002602001604051908101604052809291908181526020018280548015610b6957602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610b4b575b50505050509050919050565b815481835581811115610b9957600083815260209020610b99918101908301610b9e565b505050565b61033e91905b80821115610bb85760008155600101610ba4565b50905600a165627a7a723058206dd2e368d6f0c7701b45b4d92495e1edfa972b0b9e6ad7e7a11b0f4d9c1f03a00029`

func TestDeployContract(t *testing.T) {
	clientOnTestnet, err := cpchain.NewCPChain(cpchain.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	wallet := clientOnTestnet.LoadWallet(keystorePath)

	client, err := cpcclient.Dial(endpoint)

	if err != nil {
		t.Fatal(err)
	}

	Key, err := wallet.GetKey(password)
	if err != nil {
		t.Fatal(err)
	}

	auth := contract.NewTransactor(Key.PrivateKey, new(big.Int).SetUint64(40))
	_, _, _, err = cpchain.DeployContract2(Abi, Bin, auth, client)
	if err != nil {
		t.Fatal(err)
	}
}
