package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSyncCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		_ = cmd.Context()
		fmt.Println("not implemented yet")
	}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync all your posts",
		Args:  cobra.ExactArgs(0),
		Run:   run,
	}
	return cmd
}
