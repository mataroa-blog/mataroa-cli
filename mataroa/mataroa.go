package mataroa

import (
	"log"
	"net/http"
	"time"

	"git.sr.ht/~glorifiedgluer/mata/config"
)

type Client struct {
	endpoint string
	key      string
	HTTP     *http.Client
}

func NewMataroaClient() *Client {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{
		Timeout: time.Minute,
	}
	return &Client{
		endpoint: config.Endpoint,
		key:      config.Key,
		HTTP:     httpClient,
	}
}
