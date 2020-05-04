package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config holds options for the tool
type Config struct {
	ProjectName           string `json:"projectName"`
	DockerComposeFilePath string `json:"dockerComposeFilePath"`
	GoVersion             string `json:"goVersion"`
	GoDockerfileTemplate  string `json:"goDockerfileTemplate"`
	DatabaseService       string `json:"databaseService"`
}

var c = &Config{}

// Init reads the config file and loads it into a global variable
func Init() error {
	b, err := ioutil.ReadFile("./tools/bolt/config.json")
	if err != nil {
		return fmt.Errorf("failed to read config.json: %w", err)
	}

	if err := json.Unmarshal(b, c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Get returns the initialised config struct
func Get() *Config {
	if c == nil {
		panic("config not loaded")
	}

	return c
}
