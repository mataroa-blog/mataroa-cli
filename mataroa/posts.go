package mataroa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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
	OK bool `json:"ok,omitempty"`
}

type PostsCreateResponse struct {
	PostsBaseResponse
	Slug string `json:"slug,omitempty"`
	URL  string `json:"url,omitempty"`
}

type PostsEditResponse struct {
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

func NewMataroaRequest(ctx context.Context, client *Client, method, url string, body io.Reader) (*http.Response, error) {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return &http.Response{}, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.key))

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return resp, fmt.Errorf("error making request: %s", err)
	}

	return resp, nil
}

func CreatePost(ctx context.Context, client *Client, post PostsCreateResquest) (PostsCreateResponse, error) {
	body, err := json.Marshal(post)
	if err != nil {
		return PostsCreateResponse{}, fmt.Errorf("error marshaling post: %s", err)
	}

	resp, err := NewMataroaRequest(
		ctx,
		client,
		"POST",
		fmt.Sprintf("%s/posts", client.endpoint),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return PostsCreateResponse{}, fmt.Errorf("error creating request: %s", err)
	}
	defer resp.Body.Close()

	var response PostsCreateResponse
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	return response, nil
}

func ListPosts(ctx context.Context, client *Client) ([]Post, error) {
	var response PostsListResponse

	resp, err := NewMataroaRequest(ctx, client, "GET", fmt.Sprintf("%s/posts", client.endpoint), nil)
	if err != nil {
		return response.PostList, fmt.Errorf("error creating request: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response.PostList, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err)
	}

	return response.PostList, nil
}

func DeletePost(ctx context.Context, client *Client, slug string) (bool, error) {
	resp, err := NewMataroaRequest(ctx, client, "DELETE", fmt.Sprintf("%s/posts/%s", client.endpoint, slug), nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return false, fmt.Errorf("'%s' not found", slug)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %s", err)
	}

	var response PostsBaseResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, fmt.Errorf("error unmarshaling json: %s", err)
	}

	return response.OK, nil
}
