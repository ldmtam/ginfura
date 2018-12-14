package ginfura

// Web3Interface ...
type Web3Interface interface {
	GetBlockNumber() (uint64, error)
	ProtocolVersion() (string, error)
	Call(txCallObj TransactionCall, blkParam string) (string, error)
	GetGasPrice() (uint64, error)
	GetBalance(address string) (uint64, error)
	GetBlockByHash(blkHash string, showDetail bool) (Block, error)
	GetBlockTransactionCountByHash(blkHash string) (string, error)
	GetBlockTransactionCountByNumber(blkHash string) (string, error)
	GetCode(address string, blkParams string) (string, error)
	GetTransactionByBlockHashAndIndex(blkHash, txIndex string) (Transaction, error)
	GetTransactionByBlockNumberAndIndex(blkNumber, txIndex string) (Transaction, error)
	GetTransactionByHash(txHash string) (Transaction, error)
	GetTransactionCount(address, blkParams string) (string, error)
	GetTransactionReceipt(txHash string) (TransactionReceipt, error)
	GetUncleByBlockHashAndIndex(blkHash, index string) (UncleBlock, error)
	GetUncleByBlockNumberAndIndex(blkNumber, index string) (UncleBlock, error)
	GetUncleCountByBlockHash(blkHash string) (string, error)
	GetUncleCountByBlockNumber(blkNumber string) (string, error)
	SendRawTransaction(rawTx string) (string, error)
}
