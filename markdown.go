package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

var (
	ErrParsingFrontMatter = errors.New("error parsing frontmatter")
	ErrTitleIsEmpty       = errors.New("title cannot be empty")
	ISO8601Layout         = "2006-01-02"
)

type Frontmatter struct {
	Title string `yaml:"title"`
	// if empty, the post is a draft
	PublishedAt string `yaml:"date"`
	// if empty, the post was not published yet
	Slug string `yaml:"slug"`
}

func ParseFrontmatter(content string) (Frontmatter, error) {
	var matter Frontmatter
	_, err := frontmatter.Parse(strings.NewReader(content), &matter)
	if err != nil {
		return matter, fmt.Errorf("%w: %w", ErrParsingFrontMatter, err)
	}

	if matter.Title == "" {
		return matter, fmt.Errorf("error parsing frontmatter: %w", ErrTitleIsEmpty)
	}

	return matter, nil
}

func MarkdownToPost(content []byte, mustHaveSlug bool) (Post, error) {
	var post Post
	var metadata Frontmatter

	body, err := frontmatter.Parse(strings.NewReader(string(content)), &metadata)
	if err != nil {
		return post, fmt.Errorf("error parsing markdown file frontmatter: %s", err)
	}

	if metadata.Title == "" {
		return post, fmt.Errorf("post missing 'title' attribute")
	}

	if metadata.Slug == "" && mustHaveSlug {
		return post, fmt.Errorf("post missing 'slug' attribute")
	}

	var publishedAt string
	if metadata.PublishedAt != "" {
		t, err := time.Parse(ISO8601Layout, metadata.PublishedAt)
		if err != nil {
			return post, fmt.Errorf("post contains invalid date format '%s' should be in '%s' format",
				metadata.PublishedAt,
				ISO8601Layout,
			)
		}
		publishedAt = t.Format(ISO8601Layout)
	} else {
		publishedAt = ""
	}

	return Post{
		Body:        string(body),
		PublishedAt: publishedAt,
		Slug:        metadata.Slug,
		Title:       metadata.Title,
	}, nil
}
