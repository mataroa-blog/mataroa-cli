package mataroa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

var ErrNotFound = errors.New("not found")

type Models struct {
	Posts PostModel
}

func NewModels(client *Client) Models {
	return Models{
		Posts: PostModel{client},
	}
}

type PostModel struct {
	client *Client
}

type Post struct {
	Title       string `json:"title,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Body        string `json:"body,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	URL         string `json:"url,omitempty"`
}

// All method will return a list with all the posts
func (pm PostModel) All(ctx context.Context) (PostsListResponse, error) {
	var response PostsListResponse

	resp, err := pm.client.newMataroaRequest(ctx, "GET", "posts", nil)
	if err != nil {
		return response, fmt.Errorf("error listing posts: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	return response, nil
}

// Get method will retrieve an existing post.
func (pm PostModel) Get(ctx context.Context, slug string) (PostsGetResponse, error) {
	var response PostsGetResponse

	resp, err := pm.client.newMataroaRequest(ctx, "GET", fmt.Sprintf("posts/%s", slug), nil)
	if err != nil {
		return response, fmt.Errorf("error fetching post: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return response, ErrNotFound
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	return response, nil
}

// Delete method will delete an existing post by its slug.
func (pm PostModel) Delete(ctx context.Context, slug string) (PostsDeleteResponse, error) {
	var response PostsDeleteResponse

	resp, err := pm.client.newMataroaRequest(ctx, "DELETE", fmt.Sprintf("posts/%s", slug), nil)
	if err != nil {
		return response, fmt.Errorf("error deleting post: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return response, ErrNotFound
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling json: %s", err)
	}

	return response, nil
}

// Create method will create a new post.
func (pm PostModel) Create(ctx context.Context, post Post) (PostsCreateResponse, error) {
	var response PostsCreateResponse

	body, err := json.Marshal(post)
	if err != nil {
		return response, fmt.Errorf("error marshaling post: %s", err)
	}

	resp, err := pm.client.newMataroaRequest(ctx, "POST", "posts", bytes.NewBuffer(body))
	if err != nil {
		return response, fmt.Errorf("error creating post: %s", err)
	}
	defer resp.Body.Close()

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

// Update method will update an existing post with the new content provided
func (pm PostModel) Update(ctx context.Context, slug string, post Post) (PostsUpdateResponse, error) {
	var response PostsUpdateResponse

	body, err := json.Marshal(post)
	if err != nil {
		return response, fmt.Errorf("error updating post: %s", err)
	}

	resp, err := pm.client.newMataroaRequest(ctx, "PATCH", fmt.Sprintf("posts/%s", slug), bytes.NewBuffer(body))
	if err != nil {
		return response, fmt.Errorf("error updating post: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return response, ErrNotFound
	}

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
