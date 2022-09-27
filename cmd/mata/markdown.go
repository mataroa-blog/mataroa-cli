package main

import (
	"fmt"
	"strings"
	"time"

	"git.sr.ht/~glorifiedgluer/mata/mataroa"
	"github.com/adrg/frontmatter"
)

var ISO8601Layout = "2006-01-02"

type postFrontmatter struct {
	Title       string `yaml:"title"`
	Slug        string `yaml:"slug"`
	PublishedAt string `yaml:"published_at"`
}

func NewMarkdownToPost(content []byte) (mataroa.Post, error) {
	var post mataroa.Post
	var metadata postFrontmatter

	body, err := frontmatter.Parse(strings.NewReader(string(content)), &metadata)
	if err != nil {
		return post, fmt.Errorf("error parsing markdown file frontmatter: %s", err)
	}

	if metadata.Title == "" {
		return post, fmt.Errorf("post missing 'title' attribute")
	}

	if metadata.Slug == "" {
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

	return mataroa.Post{
		Body:        string(body),
		PublishedAt: publishedAt,
		Slug:        metadata.Slug,
		Title:       metadata.Title,
	}, nil
}

func NewPostToMarkdown(post mataroa.Post) string {
	return fmt.Sprintf(`---
title: "%s"
slug: "%s"
published_at: "%s"
---
%s`,
		post.Title,
		post.Slug,
		post.PublishedAt,
		post.Body,
	)
}

func HasPostChanged(old, new mataroa.Post) bool {
	return old.Body == new.Body &&
		old.PublishedAt == new.PublishedAt &&
		old.Slug == new.Slug &&
		old.Title == new.Title
}
