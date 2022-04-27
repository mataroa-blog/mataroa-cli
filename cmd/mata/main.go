package main

import (
	"context"
	"os"

	"git.sr.ht/~glorifiedgluer/mata/commands"
)

func main() {
	ctx := context.Background()

	cmd := commands.SetupCommands()

	if err := cmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
