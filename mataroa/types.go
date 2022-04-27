package mataroa

type PostFrontmatter struct {
	Title       string `yaml:"title"`
	Slug        string `yaml:"slug"`
	PublishedAt string `yaml:"published_at"`
}

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
