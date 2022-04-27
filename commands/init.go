package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"git.sr.ht/~glorifiedgluer/mata/config"
	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		_ = cmd.Context()

		filePath, err := xdg.ConfigFile(config.ConfigPath)
		if err != nil {
			log.Fatalf("%s: error initializing mata: %s", cmd.Use, err)
		}

		if _, err := os.Stat(filePath); err == nil {
			log.Fatalf("%s: error initializing mata: config.json already exists", cmd.Use)
		} else if errors.Is(err, os.ErrNotExist) {
			body, err := json.MarshalIndent(config.Config{
				Endpoint: "https://mataroa.blog/api",
				Key:      "your-api-key-here",
			}, "", "  ")
			if err != nil {
				log.Fatalf("%s: error initializing mata: couldn't marshal json file", cmd.Use)
			}

			err = ioutil.WriteFile(filePath, body, os.FileMode((0600)))
			if err != nil {
				log.Fatalf("%s: error initializing mata: %s", cmd.Use, err)
			}

			fmt.Printf("%s: mata initialized successfully: '%s' file created\n", cmd.Use, filePath)
		}
	}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize mata",
		Args:  cobra.ExactArgs(0),
		Run:   run,
	}
	return cmd
}
