package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func (app *application) commandsInit(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		_ = cmd.Context()

		filePath, err := xdg.ConfigFile(configPath)
		if err != nil {
			log.Fatalf("error initializing mata: %s", err)
		}

		if _, err := os.Stat(filePath); err == nil {
			log.Fatalf("error initializing mata: config.json already exists")
		} else if errors.Is(err, os.ErrNotExist) {

			body, err := json.MarshalIndent(config{
				BaseUrl: "https://mataroa.blog/api",
				Key:     "your-api-key-here",
			}, "", "  ")
			if err != nil {
				log.Fatalf("error initializing mata: couldn't marshal json file: %s", err)
			}

			err = ioutil.WriteFile(filePath, body, os.FileMode((0o600)))
			if err != nil {
				log.Fatalf("error initializing mata: %s", err)
			}

			fmt.Printf("mata initialized successfully: '%s' file created\n", filePath)
		}
	}

	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize mata",
		Run:     run,
	}

	return cmd
}

func (app *application) commandsPosts(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "posts",
		Aliases:           []string{"p"},
		Short:             "Manage posts",
		PersistentPreRunE: app.loadConfigurationPreRunE,
	}

	cmd.AddCommand(app.commandsPostsCreate(ctx))
	cmd.AddCommand(app.commandsPostsDelete(ctx))
	cmd.AddCommand(app.commandsPostsGet(ctx))
	cmd.AddCommand(app.commandsPostsList(ctx))
	cmd.AddCommand(app.commandsPostsUpdate(ctx))
	cmd.AddCommand(app.commandsPostsSync(ctx))

	return cmd
}

func (app *application) commandsPostsCreate(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%s: error creating new post: file '%s' not found\n", cmd.Use, filePath)
		}

		f, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error reading markdown file: %s", err)
		}

		post, err := NewMarkdownToPost(f)
		if err != nil {
			log.Fatalf("error creating new post: %s\n", err)
		}

		response, err := app.models.Posts.Create(ctx, post)
		if err != nil {
			log.Fatalf("error creating new post: %s\n", err)
		}

		if !response.OK {
			log.Fatalf("error creating new post: %s\n", response.Error)
		}

		fmt.Printf("'%s' created successfully: %s\n", response.Slug, response.URL)
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a post",
		Args:  cobra.ExactArgs(1),
		Run:   run,
	}

	return cmd
}

func (app *application) commandsPostsGet(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		slug := args[0]

		response, err := app.models.Posts.Get(ctx, slug)
		if err != nil {
			log.Fatalf("error getting post '%s': %s", slug, err)
		}

		if !response.OK {
			log.Fatalf("error getting post '%s': %s", slug, response.Error)
		}

		output, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Fatalf("error marshaling json: %s", err)
		}

		fmt.Println(string(output))
	}

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get a post",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}
	return cmd
}

func (app *application) commandsPostsList(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		response, err := app.models.Posts.All(ctx)
		if err != nil {
			log.Fatalf("error listing posts: %s", err)
		}

		output, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Fatalf("error marshaling json: %s", err)
		}

		fmt.Println(string(output))
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List posts",
		Args:    cobra.ExactArgs(0),
		Run:     run,
	}

	return cmd
}

func (app *application) commandsPostsDelete(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		slug := args[0]

		response, err := app.models.Posts.Delete(ctx, slug)
		if err != nil {
			log.Fatalf("error deleting post '%s': %s", slug, err)
		}

		if !response.OK {
			log.Fatalf("error deleting post '%s': %s", slug, response.Error)
		}

		fmt.Printf("post '%s' deleted successfully\n", slug)
	}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a post",
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}

	return cmd
}

func (app *application) commandsPostsUpdate(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		slug := args[0]
		filePath := args[1]

		f, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error reading markdown file: %s", err)
		}

		post, err := NewMarkdownToPost(f)
		if err != nil {
			log.Fatalf("error parsing markdown file: %s", err)
		}

		response, err := app.models.Posts.Update(ctx, slug, post)
		if err != nil {
			log.Fatalf("error updating post '%s': %s", slug, err)
		}

		if !response.OK {
			log.Fatalf("error updating post '%s': %s", slug, err)
		}

		fmt.Printf("post '%s' updated sucessfully\n", slug)
	}

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update a post",
		Aliases: []string{"u"},
		Args:    cobra.ExactArgs(2), // TODO add flags like --filename --slug
		Run:     run,
	}

	return cmd
}

func (app *application) commandsPostsSync(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		source := args[0]

		posts, err := app.models.Posts.All(ctx)
		if err != nil {
			log.Fatalf("error getting all posts: %s", err)
		}

		var postsSlugs []string
		for _, post := range posts.PostList {
			postsSlugs = append(postsSlugs, post.Slug)
		}

		var matches []string
		err = filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			ext := filepath.Ext(d.Name())
			if ext == ".md" || ext == ".markdown" {
				matches = append(matches, path)
			}

			return nil
		})
		if err != nil {
			log.Fatalf("error walking directory: %s", err)
		}

		if len(matches) == 0 {
			log.Println("no markdown files have been found")
			return
		}

		for _, filename := range matches {
			f, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Printf("error reading file '%s': %s", filename, err)
				continue
			}

			post, err := NewMarkdownToPost(f)
			if err != nil {
				log.Printf("error parsing markdown file '%s': %s", filename, err)
				continue
			}

			if ok := slices.Contains(postsSlugs, post.Slug); ok {
				response, err := app.models.Posts.Update(ctx, post.Slug, post)
				if err != nil {
					log.Printf("error updating post '%s' on filename '%s': %s", post.Slug, filename, err)
					continue
				}

				if !response.OK {
					log.Printf("error updating post '%s' on filename '%s': %s", post.Slug, filename, response.Error)
					continue
				}

				fmt.Printf("post '%s' on filename '%s' updated successfully!\n", response.Slug, filename)
				continue
			}

			response, err := app.models.Posts.Create(ctx, post)
			if err != nil {
				log.Printf("error creating post '%s' on filename '%s': %s", post.Slug, filename, err)
				continue
			}

			if !response.OK {
				log.Printf("error creating post '%s' on filename '%s': %s", post.Slug, filename, response.Error)
				continue
			}

			fmt.Printf("post '%s' on filename '%s' created successfully!\n", response.Slug, filename)
		}
	}

	cmd := &cobra.Command{
		Use:     "sync [DIRECTORY]",
		Short:   "sync all posts",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}

	return cmd
}
