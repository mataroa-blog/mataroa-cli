package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

func CommandsRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "mataroa-cli",
		Short:             "mataroa-cli is a CLI tool for mataroa.blog",
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(CommandInit())
	cmd.AddCommand(CommandsPosts())

	return cmd
}

func CommandInit() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		_ = cmd.Context()

		filePath, err := xdg.ConfigFile(ConfigurationFilePath)
		if err != nil {
			log.Fatalf("error initializing mataroa-cli: %s", err)
		}

		if _, err := os.Stat(filePath); err == nil {
			log.Fatalf("error initializing mataroa-cli: config.json already exists")
		} else if errors.Is(err, os.ErrNotExist) {

			body, err := json.MarshalIndent(Config{
				ApiUrl: "https://mataroa.blog/api",
				Token:  "your-api-key-here",
			}, "", "  ")
			if err != nil {
				log.Fatalf("error initializing mataroa-cli: couldn't marshal json file: %s", err)
			}

			err = os.WriteFile(filePath, body, os.FileMode((0o600)))
			if err != nil {
				log.Fatalf("error initializing mataroa-cli: %s", err)
			}

			fmt.Printf("mataroa-cli initialized successfully: '%s' file created\n", filePath)
		}
	}

	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize mataroa-cli",
		Run:     run,
	}

	return cmd
}

func CommandsPosts() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "posts",
		Aliases: []string{"p"},
		Short:   "Manage posts",
	}

	config, err := LoadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	client, err := NewClient(config.ApiUrl, config.Token)
	if err != nil {
		log.Fatal(err)
	}

	cmd.AddCommand(CommandPostsList(client))
	cmd.AddCommand(CommandPostsGet(client))
	cmd.AddCommand(CommandPostsDelete(client))
	cmd.AddCommand(CommandPostsUpdate(client))
	cmd.AddCommand(CommandPostsCreate(client))

	return cmd
}

func CommandPostsGet(client *Client) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		slug := args[0]

		response, err := client.Get(cmd.Context(), slug)
		if err != nil {
			log.Fatal(err)
		}

		output, err := json.MarshalIndent(response.Post, "", "  ")
		if err != nil {
			log.Fatalf("error marshaling json: %s", err)
		}

		fmt.Println(string(output))
	}

	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "get a post",
		Run:     run,
	}

	return cmd
}

func CommandPostsUpdate(client *Client) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		slug := args[0]
		filePath := args[1]

		f, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error reading markdown file: %s", err)
		}

		post, err := MarkdownToPost(f, false)
		if err != nil {
			log.Fatalf("error parsing markdown file: %s", err)
		}

		response, err := client.Update(ctx, slug, post)
		if err != nil {
			log.Fatalf("error updating post '%s': %s", slug, err)
		}

		if !response.OK {
			log.Fatalf("error updating post '%s': %s", slug, err)
		}

		fmt.Printf("updated post '%s'\n", slug)
	}

	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "update a post",
		Run:     run,
	}

	return cmd
}

func CommandPostsCreate(client *Client) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		filePath := args[0]

		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%s: error creating new post: file '%s' not found\n", cmd.Use, filePath)
		}

		f, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error reading markdown file: %s", err)
		}

		post, err := MarkdownToPost(f, false)
		if err != nil {
			log.Fatalf("error creating new post: %s\n", err)
		}

		response, err := client.Create(ctx, post)
		if err != nil {
			log.Fatalf("error creating new post: %s\n", err)
		}

		fmt.Printf("created post '%s': %s\n", response.Slug, response.URL)
	}

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "create a post",
		Run:     run,
	}

	return cmd
}

func CommandPostsDelete(client *Client) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		slug := args[0]

		_, err := client.Delete(cmd.Context(), slug)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("deleted post '%s'\n", slug)
	}

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   "delete a post",
		Run:     run,
	}

	return cmd
}

func CommandPostsList(client *Client) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		response, err := client.ListAll(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		output, err := json.MarshalIndent(response.PostList, "", "  ")
		if err != nil {
			log.Fatalf("error marshaling json: %s", err)
		}

		fmt.Println(string(output))
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List posts",
		Run:     run,
	}

	return cmd
}
