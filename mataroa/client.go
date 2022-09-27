package mataroa

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrEmptyToken    = fmt.Errorf("please, provide a non-empty token")
	ErrEmptyEndpoint = fmt.Errorf("please, provide a non-empty endpoint")
)

type Client struct {
	baseUrl    string
	token      string
	httpClient *http.Client
}

func (mc *Client) newMataroaRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", mc.baseUrl, url), body)
	if err != nil {
		return &http.Response{}, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mc.token))

	response, err := mc.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("error making request: %s", err)
	}

	return response, nil
}

// ClientBuilder interface is to abstract the builder pattern from Mataroa
// Client, gives control to the users to mock the values.
type ClientBuilder interface {
	Build() (Client, error)
	BaseUrl(string) ClientBuilder
	HttpClient(*http.Client) ClientBuilder
	Token(string) ClientBuilder
}

type clientBuilder struct {
	token      string
	endpoint   string
	httpClient *http.Client
	err        error
}

var _ ClientBuilder = (*clientBuilder)(nil)

// New method creates a new client builder
func New() ClientBuilder {
	return &clientBuilder{}
}

// Token method sets the token of the client builder
func (mb *clientBuilder) Token(token string) ClientBuilder {
	if mb.err != nil {
		return mb
	}

	if token == "" {
		mb.err = ErrEmptyToken
	}

	mb.token = token
	return mb
}

// BaseUrl method sets the base URL of the client builder
func (mb *clientBuilder) BaseUrl(endpoint string) ClientBuilder {
	if mb.err != nil {
		return mb
	}

	if endpoint == "" {
		mb.err = ErrEmptyEndpoint
	}

	mb.endpoint = endpoint
	return mb
}

// HttpClient set the http client to be used within the client calls. In case
// the httpClient argument is equal to nil, the method will generate a default
// *http.Client with sane defaults.
func (mb *clientBuilder) HttpClient(httpClient *http.Client) ClientBuilder {
	if mb.err != nil {
		return mb
	}

	mb.httpClient = httpClient
	return mb
}

// Build method returns the client with the given options
func (mb *clientBuilder) Build() (Client, error) {
	if mb.err != nil {
		return Client{}, mb.err
	}

	if mb.endpoint == "" {
		return Client{}, ErrEmptyEndpoint
	}

	if mb.token == "" {
		return Client{}, ErrEmptyToken
	}

	if mb.httpClient == nil {
		mb.httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	return Client{
		baseUrl:    mb.endpoint,
		token:      mb.token,
		httpClient: mb.httpClient,
	}, nil
}
