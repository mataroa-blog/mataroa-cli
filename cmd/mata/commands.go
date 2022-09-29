package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
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

	return cmd
}

func (app *application) commandsPostsCreate(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

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
		Use:   "create [FILENAME]",
		Short: "Create a post",
		Args:  cobra.ExactArgs(1),
		Run:   run,
	}

	return cmd
}

func (app *application) commandsPostsGet(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

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
		Use:     "get [SLUG]",
		Short:   "Get a post",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}

	return cmd
}

func (app *application) commandsPostsList(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

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
		ctx := cmd.Context()

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
		Use:     "delete [SLUG]",
		Short:   "Delete a post",
		Aliases: []string{"d"},
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}

	return cmd
}

func (app *application) commandsPostsUpdate(ctx context.Context) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		filename := args[0]

		f, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("error reading markdown file: %s", err)
		}

		post, err := NewMarkdownToPost(f)
		if err != nil {
			log.Fatalf("error parsing markdown file: %s", err)
		}

		if post.Slug == "" {
			log.Fatalf("post should have a 'slug' attribute defined")
		}

		response, err := app.models.Posts.Update(ctx, post.Slug, post)
		if err != nil {
			log.Fatalf("error updating post '%s': %s", post.Slug, err)
		}

		if !response.OK {
			log.Fatalf("error updating post '%s': %s", post.Slug, err)
		}

		fmt.Printf("post '%s' updated sucessfully\n", post.Slug)
	}

	cmd := &cobra.Command{
		Use:     "update [FILENAME]",
		Short:   "Update a post",
		Aliases: []string{"u"},
		Args:    cobra.ExactArgs(1),
	}

	cmd.Run = run

	return cmd
}
