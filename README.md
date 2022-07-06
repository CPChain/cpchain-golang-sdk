# CPChain Golang SDK
This is a light SDK for developing on the CPChain mainnet.

This project is only support get blocks and events now. If you want to develop smart contract on CPChain, please view [cpchain-cli](https://github.com/cpchain/cpchain-cli).


## OverView



## Install
before use this sdk in you code, you should install this package by go mod
```bash
go get github.com/CPChain/cpchain-golang-sdk/cpchain
```


## What can this sdk do and how to use it
* Create a cpchain instance according to the network, if only have the endpoint, it also can use the function ```GetNetWork``` to get network through endpoint
    ```go
    network, err := cpchain.GetNetWork(endpoint)

    cpchainclient, err := cpchain.NewCpchain(network)
    ```
* Get a wallet instance by importing keystorefile and password
    ```go
    cpchainclient, err := cpchain.NewCpchain(network)

    wallet, err := cpchainclient.LoadWallet(keystoreFilePath, password) // keystoreFilePath: Where the keystore file for your account is stored
    ```
* Tranfer amount to target address through a wallet instance
    ```go
    wallet, err := cpchainclient.LoadWallet(keystoreFilePath, password)

    tx, err := wallet.Transfer(targetAddress, value)
    ```
* Create a new account on the chain, return address and generate a keystore file in the path(dirpath) where you want to store it
    ```go
    cpchainclient, err := cpchain.NewCpchain(network)
    
    address, err := cpchainclient.CreateAccount(path, password)
    ```
* Deploy Contract according to the contract abi and bin ,at the same time, you need a wallet instance to send this transaction
    ```go
    cpchainclient, err := cpchain.NewCpchain(network)

    address, tx, err := cpchainclient.DepolyContract(abi, bin, wallet)

    address, tx, err := cpchainclient.DepolyContractByFile(path, wallet) //you also can deploy contract through the contract.json that build by solidity
    ```
* 
## Use cli (example)

### Create account

```bash
go run cmd/account/account.go new -keystorepath ./fixture/account
```

### Transfer
```bash
go run cmd/transfer/transfer.go -endpoint https://civilian.testnet.cpchain.io -keystore ./fixtures/keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977 -to a565060b9f2990262709075614ecec479ddf2bc7 -value 1
```
### DeployContract
```bash
go run cmd/contract/contract.go -endpoint https://civilian.testnet.cpchain.io -keystore ./fixtures/keystore/UTC--2022-06-09T05-48-04.258507200Z--52c5323efb54b8a426e84e4b383b41dcb9f7e977 -contractfile ./fixtures/contract/helloworld.json
```