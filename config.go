package main

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	Collection string   `json:"collection"`
	Tags       []string `json:"tags"`
	Resolution string   `json:"resolution"`
	User       string   `json:"user"`
}

func (config *Config) getUrl() string {
	url := "https://source.unsplash.com"

	if config.User != "" {
		url += "/user/" + config.User
	} else if config.Collection != "" {
		url += "/collection/" + config.Collection
	} else {
		url += "/random"
	}

	if config.Resolution != "" {
		url += "/" + config.Resolution
	}

	if len(config.Tags[:]) > 0 {
		url += "?" + strings.Join(config.Tags[:], ",")
	}

	return url
}

func LoadConfig(filePath string) (Config, error) {
	var config Config

	configFile, fileErr := os.Open(filePath)
	if fileErr != nil {
		return config, fileErr
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)

	decodeErr := jsonParser.Decode(&config)
	if decodeErr != nil {
		return config, decodeErr
	}

	return config, nil
}
