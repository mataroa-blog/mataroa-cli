package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrg/xdg"
)

const (
	ConfigurationFilePath = "mataroa-cli/config.json"
	EnvVariableApiUrl     = "MATAROA_API_URL"
	EnvVariableToken      = "MATAROA_TOKEN"
)

type Config struct {
	ApiUrl string `json:"api_url"`
	Token  string `json:"token"`
}

func LoadConfiguration() (Config, error) {
	var config Config

	filePath, err := xdg.ConfigFile(ConfigurationFilePath)
	if err != nil {
		return config, fmt.Errorf("error finding config file: %s", err)
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %s", err)
	}

	err = json.Unmarshal(f, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshaling config file: %s", err)
	}

	if config.ApiUrl == "" {
		config.ApiUrl = "https://mataroa.blog/api"
	}

	if config.Token == "" {
		return config, fmt.Errorf(`'key' cannot be empty on '%s'`, ConfigurationFilePath)
	}

	// environment variables will override the configuration file values
	if os.Getenv(EnvVariableApiUrl) != "" {
		config.ApiUrl = os.Getenv(EnvVariableApiUrl)
	}

	if os.Getenv(EnvVariableToken) != "" {
		config.Token = os.Getenv(EnvVariableToken)
	}

	return config, nil
}
