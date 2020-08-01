package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config holds options for the tool
type Config struct {
	Database              Database            `json:"database"`
	ProjectName           string              `json:"projectName"`
	DockerComposeFilePath string              `json:"dockerComposeFilePath"`
	Groups                map[string][]string `json:"groups"`
}

// Database holds config for the database service
type Database struct {
	Service          string `json:"service"`
	AdminService     string `json:"adminService"`
	AdminServicePath string `json:"adminServicePath"`
	Engine           string `json:"engine"`
	Username         string `json:"username"`
	Password         string `json:"password"`
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
