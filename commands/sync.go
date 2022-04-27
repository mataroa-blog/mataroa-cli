package commands

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"

	"git.sr.ht/~glorifiedgluer/mata/mataroa"
	"github.com/spf13/cobra"
)

func isMarkdownFile(path string) bool {
	return filepath.Ext(path) == ".md" ||
		filepath.Ext(path) == ".markdown"
}

func newSyncCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		c := mataroa.NewMataroaClient()

		if err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			isMd := !d.IsDir() && isMarkdownFile(path)
			if !isMd {
				return nil
			}

			f, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("%s: error reading markdown file '%s': %s", cmd.Use, path, err)
			}

			post, err := mataroa.NewPost(f)
			if err != nil {
				log.Fatalf("%s: error parsing post '%s': %s\n", cmd.Use, path, err)
			}

			_, err = c.PostBySlug(ctx, post.Slug)
			if err != nil {
				log.Printf("%s: couldn't find any post with the slug '%s'", cmd.Use, post.Slug)
				return nil
			}

			resp, err := c.UpdatePost(ctx, post.Slug, post)
			if err != nil {
				log.Fatalf("%s: error creating new post '%s': %s", cmd.Use, path, err)
			}

			if !resp.OK {
				log.Printf("%s: unable to update post '%s': %s", cmd.Use, path, resp.Error)
				return nil
			}

			fmt.Printf("%s: '%s' updated post successfully!\n", cmd.Use, resp.Slug)
			fmt.Printf("%s\n", resp.URL)
			fmt.Printf("\n")

			return nil
		}); err != nil {
			log.Fatalf("%s: %s", cmd.Use, err)
		}
	}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync all your posts",
		Args:  cobra.ExactArgs(0),
		Run:   run,
	}
	return cmd
}
