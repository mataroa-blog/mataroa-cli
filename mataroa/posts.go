package mataroa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (mc *Client) CreatePost(ctx context.Context, post PostsCreateResquest) (PostsCreateResponse, error) {
	body, err := json.Marshal(post)
	if err != nil {
		return PostsCreateResponse{}, fmt.Errorf("error marshaling post: %s", err)
	}

	resp, err := mc.newMataroaRequest(ctx, "POST", "posts", bytes.NewBuffer(body))
	if err != nil {
		return PostsCreateResponse{}, fmt.Errorf("error creating post: %s", err)
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

func (mc *Client) ListPosts(ctx context.Context) ([]Post, error) {
	var response PostsListResponse

	resp, err := mc.newMataroaRequest(ctx, "GET", "posts", nil)
	if err != nil {
		return response.PostList, fmt.Errorf("error listing posts: %s", err)
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

func (mc *Client) DeletePost(ctx context.Context, slug string) (bool, error) {
	resp, err := mc.newMataroaRequest(ctx, "DELETE", fmt.Sprintf("posts/%s", slug), nil)
	if err != nil {
		return false, fmt.Errorf("error deleting post: %s", err)
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

func (mc *Client) UpdatePost(ctx context.Context, slug string, post Post) (bool, error) {
	body, err := json.Marshal(post)
	if err != nil {
		return false, fmt.Errorf("error updating post: %s", err)
	}

	resp, err := mc.newMataroaRequest(ctx, "PATCH", fmt.Sprintf("posts/%s", slug), bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("error updating post: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return false, fmt.Errorf("'%s' not found", slug)
	}

	body, err = ioutil.ReadAll(resp.Body)
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
