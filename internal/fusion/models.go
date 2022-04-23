package fusion

import "math/big"

type (
	Dpor struct {
		Proposers  []string `json:"proposers"`
		Seal       string   `json:"seal"`
		Sigs       []string `json:"sigs"`
		Validators []string `json:"validators"`
	}
)

type HexString string

type (
	// 通过 JSON-RPC 直接获取的 Transaction
	RawTx struct {
		BlockHash        HexString `json:"blockHash"`
		BlockNumber      HexString `json:"blockNumber"`
		From             HexString `json:"from"`
		Gas              HexString `json:"gas"`
		GasPrice         HexString `json:"gasPrice"`
		Hash             HexString `json:"hash"`
		Input            HexString `json:"input"`
		Nonce            HexString `json:"nonce"`
		R                HexString `json:"r"`
		S                HexString `json:"s"`
		V                HexString `json:"v"`
		To               HexString `json:"to"`
		TransactionIndex HexString `json:"transactionIndex"`
		Type             HexString `json:"type"`
		Value            HexString `json:"value"`
	}
	// Tx
	Tx struct {
		BlockHash        HexString `json:"blockHash"`
		BlockNumber      uint64    `json:"blockNumber"`
		From             HexString `json:"from"`
		Gas              uint64    `json:"gas"`
		GasPrice         uint64    `json:"gasPrice"`
		Hash             HexString `json:"hash"`
		Input            HexString `json:"input"`
		Nonce            uint64    `json:"nonce"`
		R                HexString `json:"r"`
		S                HexString `json:"s"`
		V                HexString `json:"v"`
		To               HexString `json:"to"`
		TransactionIndex uint64    `json:"transactionIndex"`
		Type             uint64    `json:"type"`
		Value            *big.Int  `json:"value"`
	}
)

type (
	// 通过 JSON-RPC 直接获取的、fullTx 为 false，即只包含交易 Hash 的区块
	RawBlock struct {
		Dpor             Dpor        `json:"dpor"`
		Hash             HexString   `json:"hash"`
		GasLimit         HexString   `json:"gasLimit"`
		GasUsed          HexString   `json:"gasUsed"`
		ExtraData        HexString   `json:"extraData"`
		LogsBloom        HexString   `json:"logsBloom"`
		Miner            HexString   `json:"miner"`
		Number           HexString   `json:"number"`
		ParentHash       HexString   `json:"parentHash"`
		ReceiptsRoot     HexString   `json:"receiptsRoot"`
		Size             HexString   `json:"size"`
		StateRoot        HexString   `json:"stateRoot"`
		Timestamp        HexString   `json:"timestamp"`
		TransactionsRoot HexString   `json:"transactionsRoot"`
		Transactions     []HexString `json:"transactions"`
	}
	// 通过 JSON-RPC 直接获取，但指定了 fullTx，即包含了整个区块的信息
	RawBlockWithFullTxs struct {
		RawBlock
		Transactions []RawTx `json:"transactions"`
	}
	// 区块基本信息
	Block struct {
		Hash             HexString `json:"hash"`
		GasLimit         uint64    `json:"gasLimit"`
		GasUsed          uint64    `json:"gasUsed"`
		ExtraData        HexString `json:"extraData"`
		LogsBloom        HexString `json:"logsBloom"`
		Miner            string    `json:"miner"`
		Number           uint64    `json:"number"`
		ParentHash       HexString `json:"parentHash"`
		ReceiptsRoot     HexString `json:"receiptsRoot"`
		Size             uint64    `json:"size"`
		StateRoot        HexString `json:"stateRoot"`
		Timestamp        uint64    `json:"timestamp"`
		TransactionsRoot HexString `json:"transactionsRoot"`
	}
	// 包含交易的区块
	FullBlock struct {
		Block
		Transactions []Tx `json:"transactions"`
	}
)
