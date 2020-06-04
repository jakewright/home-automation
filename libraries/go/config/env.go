package config

import (
	"github.com/vrischmann/envconfig"

	"github.com/jakewright/home-automation/libraries/go/slog"
)

// Load populates the given struct with config from the environment
func Load(conf interface{}) {
	if err := envconfig.Init(conf); err != nil {
		slog.Panicf("Failed to load config from environment: %v", err)
	}
}
