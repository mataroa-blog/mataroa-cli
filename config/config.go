package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/adrg/xdg"
)

const (
	ConfigPath = "mata/config.json"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"`
}

func LoadConfig() (Config, error) {
	filePath, err := xdg.ConfigFile(ConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("error finding config file: %s", err)
	}

	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %s", err)
	}

	var config Config
	err = json.Unmarshal(f, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config file: %s", err)
	}

	if config.Endpoint == "" {
		config.Endpoint = "https://mataroa.blog/api"
	}

	if config.Key == "" {
		return Config{}, fmt.Errorf(`'key' cannot be empty on 'config.json'`)
	}

	return config, nil
}
