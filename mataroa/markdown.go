package mataroa

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

var (
	ISO8601Layout = "2006-01-02"
)

type Post struct {
	Title       string `json:"title,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Body        string `json:"body,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	URL         string `json:"url,omitempty"`
}

func NewPost(filePath string) (Post, error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Post{}, fmt.Errorf("error reading markdown file: %s", err)
	}

	var metadata PostFrontmatter
	rest, err := frontmatter.Parse(strings.NewReader(string(f)), &metadata)
	if err != nil {
		return Post{}, fmt.Errorf("error parsing markdown file frontmatter: %s", err)
	}

	if metadata.Title == "" {
		return Post{}, fmt.Errorf("post '%s' missing 'title' attribute", filePath)
	}

	if metadata.Slug == "" {
		return Post{}, fmt.Errorf("post '%s' missing 'slug' attribute", filePath)
	}

	var publishedAt string
	if metadata.PublishedAt != "" {
		t, err := time.Parse(ISO8601Layout, metadata.PublishedAt)
		if err != nil {
			return Post{}, fmt.Errorf("post '%s' contains invalid date format '%s'",
				filePath,
				metadata.PublishedAt,
			)
		}
		publishedAt = t.Format(ISO8601Layout)
	} else {
		publishedAt = ""
	}

	return Post{
		Body:        string(rest),
		PublishedAt: publishedAt,
		Slug:        metadata.Slug,
		Title:       metadata.Title,
	}, nil
}

func (p Post) ToMarkdown() string {
	return fmt.Sprintf(`---
title: %s
slug: %s
published_at: %s
---
%s
`,
		p.Title,
		p.Slug,
		p.PublishedAt,
		p.Body,
	)
}
