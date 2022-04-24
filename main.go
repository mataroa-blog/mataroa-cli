package main

import (
	"context"
	"os"

	"git.sr.ht/~glorifiedgluer/mata/commands"
	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:               "mata",
		Short:             "mata is a CLI tool for mataroa.blog",
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	}

	cmd.AddCommand(commands.NewInitCommand())
	cmd.AddCommand(commands.NewPostsCommand())
	cmd.AddCommand(commands.NewSyncCommand())

	if err := cmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
