package ginfura

import (
	"github.com/pkg/errors"
)

// common errors
var (
	errNotEthereumAddress             = errors.New("input is not an ethereum address")
	errNotOpenWebsocketConnection     = errors.New("websocket connection is not yet opened")
	errNotSubscribePendingTransaction = errors.New("pending transactions is not yet subscribed")
	errNotSubscribeNewHeads           = errors.New("new heads is not yet subscribed")
	errNotSubscribeLogs               = errors.New("logs event is not yet subscribed")
	errAlreadySubscribe               = errors.New("already subscribe the topic")
)

// subscription types
const (
	NewHeads               = "newHeads"
	Logs                   = "logs"
	NewPendingTransactions = "newPendingTransactions"
	Syncing                = "syncing"
)

// Block type of a block in ethereum blockchain
type Block struct {
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	Number           string   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             string   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	Transactions     []string `json:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
}

// Transaction ...
type Transaction struct {
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasPrice         string `json:"gasPrice"`
	Gas              string `json:"gas"`
	Input            string `json:"input"`
}

// TransactionCall ...
type TransactionCall struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value,omitempty"`
	Data     string `json:"data,omitempty"`
}

// TransactionReceipt ...
type TransactionReceipt struct {
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	From              string `json:"from"`
	To                string `json:"to"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	ContractAddress   string `json:"contractAddress"`
	Logs              []Log  `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
}

// Log ...
type Log struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

// UncleBlock ...
type UncleBlock struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Miner            string   `json:"miner"`
	Difficulty       string   `json:"difficulty"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             string   `json:"size"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Timestamp        string   `json:"timestamp"`
	Uncles           []string `json:"uncles"`
}

type blkInfoResp struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Block  `json:"result"`
}

type uncleBlkResp struct {
	JSONRPC string     `json:"jsonrpc"`
	ID      int        `json:"id"`
	Result  UncleBlock `json:"result"`
}

type txInfoResp struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  Transaction `json:"result"`
}

type txReceiptResp struct {
	JSONRPC string             `json:"jsonrpc"`
	ID      int                `json:"id"`
	Result  TransactionReceipt `json:"result"`
}

type response struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

//////////////// Websocket //////////////////
type txPendingParams struct {
	Subscription string `json:"subscription"`
	Result       string `json:"result"`
}

type txPendingResp struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  txPendingParams `json:"params"`
}

type unsubscribeResp struct {
	ID      int    `json:"id"`
	JSONRPC string `json":"jsonrpc"`
	Result  bool   `json:"result"`
}

type subscriptionResp struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}

type newHeadResult struct {
	Difficulty       string `json:"difficulty"`
	ExtraData        string `json:"extraData"`
	GasLimit         string `json:"gasLimit"`
	GasUsed          string `json:"gasUsed"`
	LogsBloom        string `json:"logsBloom"`
	Miner            string `json:"miner"`
	Nonce            string `json:"nonce"`
	Number           string `json:"number"`
	ParentHash       string `json:"parentHash"`
	ReceiptRoot      string `json:"receiptRoot"`
	Sha3Uncles       string `json:"sha3Uncles"`
	StateRoot        string `json:"stateRoot"`
	Timestamp        string `json:"timestamp"`
	TransactionsRoot string `json:"transactionsRoot"`
}
type newHeadParams struct {
	Result       newHeadResult `json:"result"`
	Subscription string        `json:"subscription"`
}

type newHeadsResp struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  newHeadParams `json:"params"`
}

type logsResult struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

type logsParams struct {
	Subscription string     `json:"subscription"`
	Result       logsResult `json:"result"`
}

type logsResp struct {
	JSONRPC string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  logsParams `json:"params"`
}

type LogRequestParams struct {
	Address []string `json:"address"`
	Topics  []string `json:"topics"`
}
