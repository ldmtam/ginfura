package ginfura

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Web3 struct {
	url    string
	client *http.Client
}

// NewWeb3Instance return new instance of ethereum api.
func NewWeb3Instance(network string, projectID string) *Web3 {
	return &Web3{
		url:    fmt.Sprintf("https://%s.infura.io/v3/%s", network, projectID),
		client: &http.Client{},
	}
}

func (e *Web3) GetBlockNumber() (uint64, error) {
	result := response{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	blkNumber, _ := strconv.ParseUint(result.Result, 0, 64)

	return blkNumber, nil
}

func (e *Web3) ProtocolVersion() (string, error) {
	result := response{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_protocolVersion",
		"params":  []interface{}{},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func validateTxCall(txCallObj TransactionCall) bool {
	if txCallObj.To == "" {
		return false
	}
	return true
}

func (e *Web3) Call(txCallObj TransactionCall, blkParam string) (string, error) {
	result := response{}

	if _, err := strconv.ParseUint(blkParam, 0, 64); err != nil && blkParam != "latest" && blkParam != "pending" && blkParam != "earliest" {
		return "", errors.New("Block param should be number or `pending`, `latest`, `earliest`")
	}

	if ok := validateTxCall(txCallObj); !ok {
		return "", errors.New("Must define `to` field")
	}

	txCallJSON, err := json.Marshal(txCallObj)
	if err != nil {
		return "", err
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_call",
		"params":  []interface{}{txCallJSON, blkParam},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetGasPrice() (uint64, error) {
	result := response{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_gasPrice",
		"params":  []interface{}{},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	gasPrice, _ := strconv.ParseUint(result.Result, 0, 64)

	return gasPrice, nil
}

func (e *Web3) GetBalance(address string) (uint64, error) {
	result := response{}

	if !isHexAddress(address) {
		return 0, errNotEthereumAddress
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  []interface{}{address, "latest"},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	balance, _ := strconv.ParseUint(result.Result, 0, 64)

	return balance, nil
}

func (e *Web3) GetBlockByHash(blkHash string, showDetail bool) (Block, error) {
	result := blkInfoResp{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByHash",
		"params":  []interface{}{blkHash, showDetail},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return Block{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetBlockTransactionCountByHash(blkHash string) (string, error) {
	result := response{}
	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockTransactionCountByHash",
		"params":  []interface{}{blkHash},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetBlockTransactionCountByNumber(blkNumber string) (string, error) {
	result := response{}

	if _, err := strconv.ParseUint(blkNumber, 0, 64); err != nil && blkNumber != "latest" && blkNumber != "pending" && blkNumber != "earliest" {
		return "", errors.New("Block param should be number or `pending`, `latest`, `earliest`")
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockTransactionCountByNumber",
		"params":  []interface{}{blkNumber},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetCode(address, blockParam string) (string, error) {
	result := response{}

	if !isHexAddress(address) {
		return "", errNotEthereumAddress
	}

	if _, err := strconv.ParseUint(blockParam, 0, 64); err != nil && blockParam != "latest" && blockParam != "pending" && blockParam != "earliest" {
		return "", errors.New("Block param should be number or `pending`, `latest`, `earliest`")
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getCode",
		"params":  []interface{}{address, blockParam},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetTransactionByBlockHashAndIndex(blkHash, index string) (Transaction, error) {
	result := txInfoResp{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionByBlockHashAndIndex",
		"params":  []interface{}{blkHash, index},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return Transaction{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetTransactionByBlockNumberAndIndex(blkNumber, index string) (Transaction, error) {
	result := txInfoResp{}

	if _, err := strconv.ParseUint(blkNumber, 0, 64); err != nil && blkNumber != "latest" && blkNumber != "pending" && blkNumber != "earliest" {
		return Transaction{}, errors.New("Block param should be number or `pending`, `latest`, `earliest`")
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionByBlockNumberAndIndex",
		"params":  []interface{}{blkNumber, index},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return Transaction{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetTransactionByHash(txHash string) (Transaction, error) {
	result := txInfoResp{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionByHash",
		"params":  []interface{}{txHash},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return Transaction{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetTransactionCount(address, blkParams string) (string, error) {
	result := response{}

	if !isHexAddress(address) {
		return "", errNotEthereumAddress
	}

	if _, err := strconv.ParseUint(blkParams, 0, 64); err != nil && blkParams != "latest" && blkParams != "pending" && blkParams != "earliest" {
		return "", errors.New("Block param should be number or `pending`, `latest`, `earliest`")
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionCount",
		"params":  []interface{}{address, blkParams},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetTransactionReceipt(txHash string) (TransactionReceipt, error) {
	result := txReceiptResp{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionReceipt",
		"params":  []interface{}{txHash},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return TransactionReceipt{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetUncleByBlockHashAndIndex(blkHash, index string) (UncleBlock, error) {
	result := uncleBlkResp{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getUncleByBlockHashAndIndex",
		"params":  []interface{}{blkHash, index},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return UncleBlock{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetUncleByBlockNumberAndIndex(blkNumber, index string) (UncleBlock, error) {
	result := uncleBlkResp{}

	if _, err := strconv.ParseUint(blkNumber, 0, 64); err != nil && blkNumber != "latest" && blkNumber != "pending" && blkNumber != "earliest" {
		return UncleBlock{}, errors.New("Block param should be number or `pending`, `latest`, `earliest`")
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getUncleByBlockNumberAndIndex",
		"params":  []interface{}{blkNumber, index},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return UncleBlock{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetUncleCountByBlockHash(blkHash string) (string, error) {
	result := response{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getUncleCountByBlockHash",
		"params":  []interface{}{blkHash},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) GetUncleCountByBlockNumber(blkNumber string) (string, error) {
	result := response{}

	if _, err := strconv.ParseUint(blkNumber, 0, 64); err != nil && blkNumber != "latest" && blkNumber != "pending" && blkNumber != "earliest" {
		return "", errors.New("Block param should be number or `pending`, `latest`, `earliest`")
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getUncleCountByBlockNumber",
		"params":  []interface{}{blkNumber},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}

func (e *Web3) SendRawTransaction(rawTx string) (string, error) {
	result := response{}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_sendRawTransaction",
		"params":  []interface{}{rawTx},
		"id":      1,
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := e.client.Post(e.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Result, nil
}
