# CPChain Golang SDK

This is a light SDK for developing on the CPChain mainnet.

This project is only support get blocks and events now. If you want to develop smart contract on CPChain, please view [cpchain-cli](https://github.com/cpchain/cpchain-cli).

## Interface

## Function 
### NewCpchain
this function will create a cpchain(interface) instance according to network
example:
``` go
client, err := cpchain.NewCpchain(cpchain.Testnet)
```




## Use cli (example)

### Create account

```bash
go run cmd/account/account.go new -keystorepath ./fixture/account
```

### Transfer

go run cmd/transfer/transfer.go -endpoint https://civilian.testnet.cpchain.io -keystore ./fixtures/keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977 -to a565060b9f2990262709075614ecec479ddf2bc7 -value 1

### DeployContract
go run cmd/contract/contract.go -endpoint https://civilian.testnet.cpchain.io -keystore ./fixtures/keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977 -contractfile ./fixtures/contract/helloworld.json