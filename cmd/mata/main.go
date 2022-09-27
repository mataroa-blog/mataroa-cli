package main

import (
	"context"
	"log"

	"git.sr.ht/~glorifiedgluer/mata/mataroa"
	"github.com/spf13/cobra"
)

type application struct {
	config config
	models mataroa.Models
}

func main() {
	app := &application{}

	ctx := context.Background()
	if err := app.run(ctx); err != nil {
		log.Fatalf("error running application: %s", err)
	}
}

func (app *application) run(ctx context.Context) error {
	cmd := &cobra.Command{
		Use:               "mata",
		Short:             "mata is a CLI tool for mataroa.blog",
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(app.commandsPosts(ctx))
	cmd.AddCommand(app.commandsInit(ctx))

	return cmd.ExecuteContext(ctx)
}
