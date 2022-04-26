package mataroa

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

func (mc *Client) newMataroaRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", mc.endpoint, url), body)
	if err != nil {
		return &http.Response{}, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mc.key))

	resp, err := mc.HTTP.Do(req)
	if err != nil {
		return resp, fmt.Errorf("error making request: %s", err)
	}

	return resp, nil
}
