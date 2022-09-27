package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"git.sr.ht/~glorifiedgluer/mata/mataroa"
	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

const (
	configPath  = "mata/config.json"
	envEndpoint = "MATAROA_ENDPOINT"
	envKey      = "MATAROA_KEY"
)

type config struct {
	BaseUrl string `json:"base_url"`
	Key     string `json:"key"`
}

func (app *application) loadConfigurationPreRunE(cmd *cobra.Command, args []string) error {
	return app.loadConfiguration()
}

func (app *application) loadConfiguration() error {
	if os.Getenv(envEndpoint) != "" && os.Getenv(envKey) != "" {
		app.config.BaseUrl = os.Getenv(envEndpoint)
		app.config.Key = os.Getenv(envKey)
		return nil
	}

	filePath, err := xdg.ConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("error finding config file: %s", err)
	}

	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %s", err)
	}

	var config config
	err = json.Unmarshal(f, &config)
	if err != nil {
		return fmt.Errorf("error unmarshaling config file: %s", err)
	}

	if config.BaseUrl == "" {
		config.BaseUrl = "https://mataroa.blog/api"
	}

	if config.Key == "" {
		return fmt.Errorf(`'key' cannot be empty on '%s'`, configPath)
	}

	client, err := mataroa.New().Token(app.config.Key).BaseUrl(app.config.BaseUrl).Build()
	if err != nil {
		log.Fatalf("error build mataroa client: %s", err)
	}

	app.models = mataroa.NewModels(&client)

	return nil
}
