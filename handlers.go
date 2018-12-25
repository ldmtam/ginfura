package ginfura

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/orcaman/concurrent-map"
)

type subscription struct {
	conn           *websocket.Conn
	subscriptionID string
}

// Ginfura ...
type Ginfura struct {
	// Http connection
	url    string
	client *http.Client

	// Websocket connection
	wsURL           string
	subscriptionMap cmap.ConcurrentMap // SubscriptionType => subscription
}

// NewGinfura return new instance of ginfura api.
func NewGinfura(network string, projectID string) *Ginfura {
	var url string
	var wsURL string
	subsMap := cmap.New()

	if projectID == "" {
		url = fmt.Sprintf("https://%s.infura.io/", network)
		wsURL = fmt.Sprintf("wss://%s.infura.io/ws", network)
	} else {
		url = fmt.Sprintf("https://%s.infura.io/v3/%s", network, projectID)
		wsURL = fmt.Sprintf("wss://%s.infura.io/v3/%s/ws", network, projectID)
	}

	return &Ginfura{
		url:             url,
		wsURL:           wsURL,
		client:          &http.Client{},
		subscriptionMap: subsMap,
	}
}
