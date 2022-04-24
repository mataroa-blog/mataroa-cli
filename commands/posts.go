package commands

import (
	"errors"
	"fmt"
	"log"
	"os"

	"git.sr.ht/~glorifiedgluer/mata/mataroa"
	"github.com/spf13/cobra"
)

func NewPostsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "posts",
		Short: "Manage posts",
	}

	cmd.AddCommand(NewPostsCreateCommand())
	cmd.AddCommand(NewPostsDeleteCommand())
	cmd.AddCommand(NewPostsEditCommand())
	cmd.AddCommand(NewPostsListCommand())

	return cmd
}

func NewPostsCreateCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		filePath := args[0]

		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("error creating new post: file '%s' not found\n", filePath)
		}

		post, err := mataroa.NewPost(filePath)
		if err != nil {
			log.Fatalf("error creating new post: %s\n", err)
		}

		c := mataroa.NewMataroaClient()

		resp, err := mataroa.CreatePost(ctx, c, mataroa.PostsCreateResquest{
			Title:       post.Title,
			PublishedAt: post.PublishedAt,
			Body:        post.Body,
		})
		if err != nil {
			log.Fatalf("error creating new post: %s\n", err)
		}

		if resp.OK {
			fmt.Printf("post created successfully! %s\n", resp.URL)
		}
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a post",
		Args:  cobra.ExactArgs(1),
		Run:   run,
	}
	return cmd
}

func NewPostsDeleteCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		slug := args[0]

		c := mataroa.NewMataroaClient()

		ok, err := mataroa.DeletePost(ctx, c, slug)
		if err != nil {
			log.Fatal(err)
		}

		if !ok {
			log.Fatalf("could not delete '%s'\n", slug)
			return
		}

		fmt.Printf("deleted post '%s'\n", slug)
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a post",
		Args:  cobra.ExactArgs(1),
		Run:   run,
	}
	return cmd
}

func NewPostsEditCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		_ = cmd.Context()
		fmt.Println("not implemented yet")
	}

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a post",
		Args:  cobra.ExactArgs(0),
		Run:   run,
	}
	return cmd
}

func NewPostsListCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		c := mataroa.NewMataroaClient()

		posts, err := mataroa.ListPosts(ctx, c)
		if err != nil {
			log.Fatal(err)
		}

		for _, post := range posts {
			fmt.Printf("%s\n", post.Slug)
			fmt.Printf("%s - %s", post.Title, post.PublishedAt)
			fmt.Printf("\n\n")
		}
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List posts",
		Args:  cobra.ExactArgs(0),
		Run:   run,
	}
	return cmd
}
