package mataroa

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewPost(t *testing.T) {
	type args struct {
		filePath string
		content  string
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
				Title:       "Foobar",
				Slug:        "foobar",
				PublishedAt: "2022-01-02",
				Body:        "FooBar and FuzBax.",
			},
			args: args{
				content: `
---
title: Foobar
slug: foobar
published_at: 2022-01-02
---
FooBar and FuzBax.`,
			},
		},
		{
			name:    "should return post successfully with missing published_at",
			wantErr: false,
			want: Post{
				Title:       "Foobar",
				Slug:        "foobar",
				PublishedAt: "",
				Body:        "FooBar and FuzBax.",
			},
			args: args{
				content: `
---
title: Foobar
slug: foobar
---
FooBar and FuzBax.`,
			},
		},
		{
			name:    "should error on missing title",
			wantErr: true,
			want:    Post{},
			args: args{
				content: `
---
slug: foobar
published_at: 2022-01-02
---
FooBar and FuzBax.`,
			},
		},
		{
			name:    "should error on missing slug",
			wantErr: true,
			want:    Post{},
			args: args{
				content: `
---
title: Foobar
published_at: 2022-01-02
---
FooBar and FuzBax.`,
			},
		},
		{
			name:    "should error on malformated published at date",
			wantErr: true,
			want:    Post{},
			args: args{
				content: `
---
title: Foobar
slug: foobar
published_at: 01-02-2006
---
FooBar and FuzBax.`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := ioutil.TempFile(t.TempDir(), "*")
			if err != nil {
				t.Errorf("%s", err)
			}
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(tt.args.content)
			if err != nil {
				t.Errorf("%s", err)
			}

			tt.args.filePath = tmpFile.Name()

			got, err := NewPost(tt.args.filePath)
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
