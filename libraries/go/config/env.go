package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/vrischmann/envconfig"

	"github.com/jakewright/home-automation/libraries/go/slog"
)

// EnvProvider loads config from environment variables
type EnvProvider struct {
	Prefix string
}

// Get returns the config value with the given key
func (e EnvProvider) Get(key string) Value {
	key = strings.ToUpper(key)
	if e.Prefix != "" {
		key = fmt.Sprintf("%s_%s", e.Prefix, key)
	}
	raw, ok := os.LookupEnv(key)
	if !ok {
		slog.Panicf("Environment variable %s not set", key)
	}
	return Value{
		key:    key,
		raw:    raw,
		exists: ok,
	}
}

// Load populates the given struct with config from the environment
func (e EnvProvider) Load(conf interface{}) {
	if err := envconfig.InitWithPrefix(conf, e.Prefix); err != nil {
		slog.Panicf("Failed to load config from environment: %v", err)
	}
}
