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

// Ginfura ...
type Ginfura struct {
	url    string
	client *http.Client
}

// NewGinfura return new instance of ginfura api.
func NewGinfura(network string, projectID string) *Ginfura {
	var url string

	if projectID == "" {
		url = fmt.Sprintf("https://%s.infura.io/", network)
	} else {
		url = fmt.Sprintf("https://%s.infura.io/v3/%s", network, projectID)
	}

	return &Ginfura{
		url:    url,
		client: &http.Client{},
	}
}

func (e *Ginfura) GetBlockNumber() (uint64, error) {
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

func (e *Ginfura) ProtocolVersion() (string, error) {
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

func (e *Ginfura) Call(txCallObj TransactionCall, blkParam string) (string, error) {
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

func (e *Ginfura) GetGasPrice() (uint64, error) {
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

func (e *Ginfura) GetBalance(address string) (uint64, error) {
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

func (e *Ginfura) GetBlockByHash(blkHash string, showDetail bool) (Block, error) {
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

func (e *Ginfura) GetBlockTransactionCountByHash(blkHash string) (string, error) {
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

func (e *Ginfura) GetBlockTransactionCountByNumber(blkNumber string) (string, error) {
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

func (e *Ginfura) GetCode(address, blockParam string) (string, error) {
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

func (e *Ginfura) GetTransactionByBlockHashAndIndex(blkHash, index string) (Transaction, error) {
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

func (e *Ginfura) GetTransactionByBlockNumberAndIndex(blkNumber, index string) (Transaction, error) {
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

func (e *Ginfura) GetTransactionByHash(txHash string) (Transaction, error) {
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

func (e *Ginfura) GetTransactionCount(address, blkParams string) (string, error) {
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

func (e *Ginfura) GetTransactionReceipt(txHash string) (TransactionReceipt, error) {
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

func (e *Ginfura) GetUncleByBlockHashAndIndex(blkHash, index string) (UncleBlock, error) {
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

func (e *Ginfura) GetUncleByBlockNumberAndIndex(blkNumber, index string) (UncleBlock, error) {
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

func (e *Ginfura) GetUncleCountByBlockHash(blkHash string) (string, error) {
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

func (e *Ginfura) GetUncleCountByBlockNumber(blkNumber string) (string, error) {
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

func (e *Ginfura) SendRawTransaction(rawTx string) (string, error) {
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
