package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidUrl   = errors.New("invalid url")
	ErrInvalidToken = errors.New("invalid token")
	ErrNotFound     = errors.New("not found")
)

type Post struct {
	Title       string `json:"title,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Body        string `json:"body,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	URL         string `json:"url,omitempty"`
}

type PostsCreateResquest struct {
	Body        string `json:"body,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	Title       string `json:"title,omitempty"`
}

type PostsBaseResponse struct {
	OK    bool   `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
}

type PostsCreateResponse struct {
	PostsBaseResponse
	Slug string `json:"slug,omitempty"`
	URL  string `json:"url,omitempty"`
}

type PostsUpdateResponse struct {
	PostsBaseResponse
	Slug string `json:"slug,omitempty"`
	URL  string `json:"url,omitempty"`
}

type PostsDeleteResponse struct {
	PostsBaseResponse
}

type PostsListResponse struct {
	PostsBaseResponse
	PostList []Post `json:"post_list,omitempty"`
}

type PostsGetResponse struct {
	PostsBaseResponse
	Post
}

type Client struct {
	baseUrl    *url.URL
	token      string
	httpClient *http.Client
}

func NewClient(baseUrl string, token string) (*Client, error) {
	if baseUrl == "" {
		return &Client{}, ErrInvalidUrl
	}
	if token == "" {
		return &Client{}, ErrInvalidToken
	}

	u, err := url.Parse(baseUrl)
	if err != nil {
		return &Client{}, ErrInvalidUrl
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	return &Client{
		baseUrl: u, token: token, httpClient: httpClient,
	}, nil
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.baseUrl, url), body)
	if err != nil {
		return &http.Response{}, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	response, err := c.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("error making request: %s", err)
	}

	return response, nil
}

func (c *Client) ListAll(ctx context.Context) (PostsListResponse, error) {
	var response PostsListResponse

	resp, err := c.newRequest(ctx, "GET", "posts", nil)
	if err != nil {
		return response, fmt.Errorf("error listing posts: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	if !response.OK {
		return response, fmt.Errorf("error: %s", response.Error)
	}

	return response, nil
}

func (c *Client) Get(ctx context.Context, slug string) (PostsGetResponse, error) {
	var response PostsGetResponse

	resp, err := c.newRequest(ctx, "GET", fmt.Sprintf("posts/%s", slug), nil)
	if err != nil {
		return response, fmt.Errorf("error fetching post: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return response, fmt.Errorf("error: %w", ErrNotFound)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	if !response.OK {
		return response, fmt.Errorf("error: %s", response.Error)
	}

	return response, nil
}

func (c *Client) Delete(ctx context.Context, slug string) (PostsDeleteResponse, error) {
	var response PostsDeleteResponse

	resp, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("posts/%s", slug), nil)
	if err != nil {
		return response, fmt.Errorf("error deleting post: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return response, ErrNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	if !response.OK {
		return response, fmt.Errorf("error: %s", response.Error)
	}

	return response, nil
}

func (c *Client) Update(ctx context.Context, slug string, post Post) (PostsUpdateResponse, error) {
	var response PostsUpdateResponse

	body, err := json.Marshal(post)
	if err != nil {
		return response, fmt.Errorf("error updating post: %s", err)
	}

	resp, err := c.newRequest(ctx, "PATCH", fmt.Sprintf("posts/%s", slug), bytes.NewBuffer(body))
	if err != nil {
		return response, fmt.Errorf("error updating post: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return response, ErrNotFound
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	if !response.OK {
		return response, fmt.Errorf("error: %s", response.Error)
	}

	return response, nil
}

func (c *Client) Create(ctx context.Context, post Post) (PostsCreateResponse, error) {
	var response PostsCreateResponse

	body, err := json.Marshal(post)
	if err != nil {
		return response, fmt.Errorf("error marshaling post: %s", err)
	}

	resp, err := c.newRequest(ctx, "POST", "posts", bytes.NewBuffer(body))
	if err != nil {
		return response, fmt.Errorf("error creating post: %s", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	if !response.OK {
		return response, fmt.Errorf("error: %s", response.Error)
	}

	return response, nil
}
