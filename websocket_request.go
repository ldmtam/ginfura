package ginfura

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// use this to make sure that we'll kill all goroutines used to listen data from infura server
// before closing the connection.
var wg sync.WaitGroup

func (g *Ginfura) Close() {
	if g.wsClient == nil {
		return
	}
	wg.Wait()
	g.wsClient.Close()
}

func (g *Ginfura) Open() error {
	c, _, err := websocket.DefaultDialer.Dial(g.wsURL, nil)
	if err != nil {
		return err
	}
	g.wsClient = c

	return nil
}

func (g *Ginfura) SubscribePendingTransaction(ctx context.Context) (<-chan string, error) {
	if g.wsClient == nil {
		return nil, errNotOpenWebsocketConnection
	}

	g.mu.Lock()
	if _, ok := g.subscriptions[NewPendingTransactions]; ok {
		return nil, errAlreadySubscribe
	}
	g.mu.Unlock()

	txPendingQueue := make(chan string, 128)

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_subscribe",
		"params":  []interface{}{"newPendingTransactions"},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	err = g.wsClient.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go listenForPendingTx(ctx, g, txPendingQueue)

	return txPendingQueue, nil
}

func (g *Ginfura) SubscribeNewHeads(ctx context.Context) (<-chan newHeadResult, error) {
	if g.wsClient == nil {
		return nil, errNotOpenWebsocketConnection
	}

	g.mu.Lock()
	if _, ok := g.subscriptions[NewHeads]; ok {
		return nil, errAlreadySubscribe
	}
	g.mu.Unlock()

	newHeadsQueue := make(chan newHeadResult, 128)

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_subscribe",
		"params":  []interface{}{"newHeads"},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	err = g.wsClient.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go listenForNewHeads(ctx, g, newHeadsQueue)

	return newHeadsQueue, nil
}

func (g *Ginfura) unsubscribePendingTransaction() (bool, error) {
	if g.wsClient == nil {
		return false, errNotOpenWebsocketConnection
	}

	g.mu.Lock()
	subscriptionID, ok := g.subscriptions[NewPendingTransactions]
	g.mu.Unlock()
	if ok == false {
		return false, errNotSubscribePendingTransaction
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_unsubscribe",
		"params":  []interface{}{subscriptionID},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return false, err
	}
	err = g.wsClient.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return false, err
	}

	resp := unsubscribeResp{}
	for {
		_, message, err := g.wsClient.ReadMessage()
		if err != nil {
			return false, err
		}
		if err = json.Unmarshal(message, &resp); err != nil {
			return false, err
		}
		if resp.Result == true {
			delete(g.subscriptions, NewPendingTransactions)
		}
		return resp.Result, nil
	}
}

func (g *Ginfura) unsubscribeNewHeads() (bool, error) {
	if g.wsClient == nil {
		return false, errNotOpenWebsocketConnection
	}

	g.mu.Lock()
	subscriptionID, ok := g.subscriptions[NewHeads]
	g.mu.Unlock()
	if ok == false {
		return false, errNotSubscribeNewHeads
	}

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_unsubscribe",
		"params":  []interface{}{subscriptionID},
		"id":      1,
	}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return false, err
	}
	err = g.wsClient.WriteMessage(websocket.TextMessage, jsonValue)
	if err != nil {
		return false, err
	}

	resp := unsubscribeResp{}
	for {
		_, message, err := g.wsClient.ReadMessage()
		if err != nil {
			return false, err
		}
		if err = json.Unmarshal(message, &resp); err != nil {
			return false, err
		}
		if resp.Result == true {
			delete(g.subscriptions, NewHeads)
		}
		return resp.Result, nil
	}
}

func listenForPendingTx(ctx context.Context, g *Ginfura, pendingTxQueue chan<- string) {
	resp := txPendingResp{}
	i := 0
	for {
		select {
		case <-ctx.Done():
			g.unsubscribePendingTransaction()
			log.Println("Unsubscribing pending transaction event")
			wg.Done()
			return
		default:
			_, message, err := g.wsClient.ReadMessage()
			if err != nil {
				log.Println("error:", err)
				return
			}
			if i == 0 {
				subResp := subscriptionResp{}
				if err = json.Unmarshal(message, &subResp); err != nil {
					log.Println("error:", err)
					return
				}

				g.mu.Lock()
				g.subscriptions[NewPendingTransactions] = subResp.Result
				g.mu.Unlock()

				i++
			} else {
				if err = json.Unmarshal(message, &resp); err != nil {
					log.Println("error:", err)
					return
				}
				pendingTxQueue <- resp.Params.Result
			}
		}

	}
}

func listenForNewHeads(ctx context.Context, g *Ginfura, newHeadsQueue chan<- newHeadResult) {
	resp := newHeadsResp{}
	i := 0
	for {
		select {
		case <-ctx.Done():
			log.Println("Unsubscribing new heads event")
			g.unsubscribeNewHeads()
			wg.Done()
			return
		default:
			_, message, err := g.wsClient.ReadMessage()
			if err != nil {
				log.Println("error:", err)
				return
			}
			if i == 0 {
				subResp := subscriptionResp{}
				if err = json.Unmarshal(message, &subResp); err != nil {
					log.Println("error:", err)
					return
				}

				g.mu.Lock()
				g.subscriptions[NewHeads] = subResp.Result
				g.mu.Unlock()

				i++
			} else {
				if err = json.Unmarshal(message, &resp); err != nil {
					log.Println("error:", err)
					return
				}
				newHeadsQueue <- resp.Params.Result
			}
		}
	}
}
