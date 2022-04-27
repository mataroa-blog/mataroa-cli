package mataroa

import (
	"reflect"
	"testing"
)

func TestNewPost(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Post
		wantErr bool
	}{
		{
			name:    "should return post successfully",
			wantErr: false,
			want: Post{
				Title:       "test: file title",
				Slug:        "file-title",
				PublishedAt: "2022-01-02",
				Body:        "FooBar and FuzBax.",
			},
			args: args{
				content: []byte(`
---
title: "test: file title"
slug: "file-title"
published_at: "2022-01-02"
---
FooBar and FuzBax.`),
			},
		},
		{
			name:    "should return post successfully with missing published_at",
			wantErr: false,
			want: Post{
				Title:       "test: file title",
				Slug:        "file-title",
				PublishedAt: "",
				Body:        "FooBar and FuzBax.",
			},
			args: args{
				content: []byte(`
---
title: "test: file title"
slug: "file-title"
---
FooBar and FuzBax.`),
			},
		},
		{
			name:    "should error on missing title",
			wantErr: true,
			want:    Post{},
			args: args{
				content: []byte(`
---
slug: foobar
published_at: 2022-01-02
---
FooBar and FuzBax.`),
			},
		},
		{
			name:    "should error on missing slug",
			wantErr: true,
			want:    Post{},
			args: args{
				content: []byte(`
---
title: Foobar
published_at: 2022-01-02
---
FooBar and FuzBax.`),
			},
		},
		{
			name:    "should error on malformated published at date",
			wantErr: true,
			want:    Post{},
			args: args{
				content: []byte(`
---
title: Foobar
slug: foobar
published_at: 01-02-2006
---
FooBar and FuzBax.`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPost(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPost() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPost() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestPost_ToMarkdown(t *testing.T) {
	type fields struct {
		Title       string
		Slug        string
		Body        string
		PublishedAt string
		URL         string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should generate post successfully with all attributes",
			fields: fields{
				Title:       "Foobar",
				Slug:        "foobar",
				Body:        "Foobar and Fuzbax.\n",
				PublishedAt: "2006-01-02",
			},
			want: `---
title: "Foobar"
slug: "foobar"
published_at: "2006-01-02"
---
Foobar and Fuzbax.
`,
		},
		{
			name: "should generate post successfully without published_at",
			fields: fields{
				Title: "Foobar",
				Slug:  "foobar",
				Body:  "Foobar and Fuzbax.\n",
			},
			want: `---
title: "Foobar"
slug: "foobar"
published_at: ""
---
Foobar and Fuzbax.
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Post{
				Title:       tt.fields.Title,
				Slug:        tt.fields.Slug,
				Body:        tt.fields.Body,
				PublishedAt: tt.fields.PublishedAt,
				URL:         tt.fields.URL,
			}
			if got := p.ToMarkdown(); got != tt.want {
				t.Errorf("Post.ToMarkdown() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
