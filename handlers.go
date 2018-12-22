package ginfura

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Ginfura ...
type Ginfura struct {
	url           string
	wsURL         string
	client        *http.Client
	wsClient      *websocket.Conn
	subscriptions map[string]string // Subscription Type => Subscription ID.
	mu            *sync.Mutex
}

// NewGinfura return new instance of ginfura api.
func NewGinfura(network string, projectID string) (*Ginfura, error) {
	var url string
	var wsURL string
	subscriptions := make(map[string]string)

	if projectID == "" {
		url = fmt.Sprintf("https://%s.infura.io/", network)
		wsURL = fmt.Sprintf("wss://%s.infura.io/ws", network)
	} else {
		url = fmt.Sprintf("https://%s.infura.io/v3/%s", network, projectID)
		wsURL = fmt.Sprintf("wss://%s.infura.io/v3/%s/ws", network, projectID)
	}

	return &Ginfura{
		url:           url,
		wsURL:         wsURL,
		client:        &http.Client{},
		wsClient:      nil,
		subscriptions: subscriptions,
		mu:            new(sync.Mutex),
	}, nil
}
