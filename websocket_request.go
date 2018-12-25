package ginfura

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

func (g *Ginfura) SubscribePendingTransaction() (<-chan string, chan struct{}, error) {

	// initialize subscription struct for pending transaction event.
	sub := &subscription{}

	// open websocket connection
	c, _, err := websocket.DefaultDialer.Dial(g.wsURL, nil)
	if err != nil {
		return nil, nil, err
	}

	sub.conn = c

	// prepare message to send to infura server
	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_subscribe",
		"params":  []interface{}{"newPendingTransactions"},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, nil, err
	}

	// send request to infura server
	err = c.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return nil, nil, err
	}

	for {
		// Read message from infura server.
		_, message, err := c.ReadMessage()
		if err != nil {
			return nil, nil, err
		}

		subResp := subscriptionResp{}
		json.Unmarshal(message, &subResp)
		sub.subscriptionID = subResp.Result

		g.subscriptionMap.Set(NewPendingTransaction, sub)
		break
	}

	done := make(chan struct{})
	pendingTxQueue := make(chan string)
	go listenForPendingTx(g, pendingTxQueue, done)

	return pendingTxQueue, done, nil
}

func (g *Ginfura) UnSubscribePendingTransaction() {

	pendingTxSub := &subscription{}
	if tmp, ok := g.subscriptionMap.Get(NewPendingTransaction); ok {
		pendingTxSub = tmp.(*subscription)
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_unsubscribe",
		"params":  []interface{}{pendingTxSub.subscriptionID},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return
	}
	err = pendingTxSub.conn.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return
	}

	resp := unsubscribeResp{}
	for {
		_, message, err := pendingTxSub.conn.ReadMessage()
		if err != nil {
			return
		}
		json.Unmarshal(message, &resp)

		if resp.Result == true {
			pendingTxSub.conn.Close()
			g.subscriptionMap.Remove(NewPendingTransaction)
		}
		return
	}
}

func listenForPendingTx(g *Ginfura, pendingTxQueue chan<- string, done chan struct{}) {
	pendingTxSub := &subscription{}
	if tmp, ok := g.subscriptionMap.Get(NewPendingTransaction); ok {
		pendingTxSub = tmp.(*subscription)
	}

	resp := txPendingResp{}
	for {
		select {
		case <-done:
			g.UnSubscribePendingTransaction()
			return
		default:
			_, message, err := pendingTxSub.conn.ReadMessage()
			if err != nil {
				return
			}

			json.Unmarshal(message, &resp)

			pendingTxQueue <- resp.Params.Result
		}
	}
}

func (g *Ginfura) SubscribeNewHead() (<-chan newHeadResult, chan struct{}, error) {

	// initialize subscription struct for pending transaction event.
	sub := &subscription{}

	// open websocket connection
	c, _, err := websocket.DefaultDialer.Dial(g.wsURL, nil)
	if err != nil {
		return nil, nil, err
	}

	sub.conn = c

	// prepare message to send to infura server
	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_subscribe",
		"params":  []interface{}{"newHeads"},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, nil, err
	}

	// send request to infura server
	err = c.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return nil, nil, err
	}

	for {
		// Read message from infura server.
		_, message, err := c.ReadMessage()
		if err != nil {
			return nil, nil, err
		}

		subResp := subscriptionResp{}
		json.Unmarshal(message, &subResp)
		sub.subscriptionID = subResp.Result

		g.subscriptionMap.Set(NewHead, sub)
		break
	}

	done := make(chan struct{})
	newHeadQueue := make(chan newHeadResult)
	go listenForNewHead(g, newHeadQueue, done)

	return newHeadQueue, done, nil
}

func (g *Ginfura) UnSubscribeNewHead() {

	newHeadSub := &subscription{}
	if tmp, ok := g.subscriptionMap.Get(NewHead); ok {
		newHeadSub = tmp.(*subscription)
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_unsubscribe",
		"params":  []interface{}{newHeadSub.subscriptionID},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return
	}
	err = newHeadSub.conn.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return
	}

	resp := unsubscribeResp{}
	for {
		_, message, err := newHeadSub.conn.ReadMessage()
		if err != nil {
			return
		}
		json.Unmarshal(message, &resp)

		if resp.Result == true {
			newHeadSub.conn.Close()
			g.subscriptionMap.Remove(NewHead)
		}
		return
	}
}

func listenForNewHead(g *Ginfura, newHeadQueue chan<- newHeadResult, done chan struct{}) {
	newHeadSub := &subscription{}
	if tmp, ok := g.subscriptionMap.Get(NewHead); ok {
		newHeadSub = tmp.(*subscription)
	}

	resp := newHeadsResp{}
	for {
		select {
		case <-done:
			g.UnSubscribeNewHead()
			return
		default:
			_, message, err := newHeadSub.conn.ReadMessage()
			if err != nil {
				return
			}

			json.Unmarshal(message, &resp)

			newHeadQueue <- resp.Params.Result
		}
	}
}

func (g *Ginfura) SubscribeNewLog(params *LogRequestParams) (<-chan logsResult, chan struct{}, error) {

	// initialize subscription struct for pending transaction event.
	sub := &subscription{}

	// open websocket connection
	c, _, err := websocket.DefaultDialer.Dial(g.wsURL, nil)
	if err != nil {
		return nil, nil, err
	}

	sub.conn = c

	// prepare message to send to infura server
	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_subscribe",
		"params":  []interface{}{"logs", params},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, nil, err
	}

	// send request to infura server
	err = c.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return nil, nil, err
	}

	for {
		// Read message from infura server.
		_, message, err := c.ReadMessage()
		if err != nil {
			return nil, nil, err
		}

		subResp := subscriptionResp{}
		json.Unmarshal(message, &subResp)
		sub.subscriptionID = subResp.Result

		g.subscriptionMap.Set(NewLog, sub)
		break
	}

	done := make(chan struct{})
	logQueue := make(chan logsResult)
	go listenForNewLog(g, logQueue, done)

	return logQueue, done, nil
}

func (g *Ginfura) UnSubscribeNewLog() {

	newLogSub := &subscription{}
	if tmp, ok := g.subscriptionMap.Get(NewLog); ok {
		newLogSub = tmp.(*subscription)
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_unsubscribe",
		"params":  []interface{}{newLogSub.subscriptionID},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return
	}
	err = newLogSub.conn.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return
	}

	resp := unsubscribeResp{}
	for {
		_, message, err := newLogSub.conn.ReadMessage()
		if err != nil {
			return
		}
		json.Unmarshal(message, &resp)

		if resp.Result == true {
			newLogSub.conn.Close()
			g.subscriptionMap.Remove(NewLog)
		}
		return
	}
}

func listenForNewLog(g *Ginfura, logQueue chan<- logsResult, done chan struct{}) {
	newLogSub := &subscription{}
	if tmp, ok := g.subscriptionMap.Get(NewLog); ok {
		newLogSub = tmp.(*subscription)
	}

	resp := logsResp{}
	for {
		select {
		case <-done:
			g.UnSubscribeNewLog()
			return
		default:
			_, message, err := newLogSub.conn.ReadMessage()
			if err != nil {
				return
			}

			json.Unmarshal(message, &resp)

			logQueue <- resp.Params.Result
		}
	}
}
