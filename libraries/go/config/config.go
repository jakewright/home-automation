package config

import (
	"github.com/jakewright/home-automation/libraries/go/slog"
)

// Provider allows reading of config values
type Provider interface {
	Get(string) Value
	Load(interface{})
}

// DefaultProvider is a global instance of a Provider
var DefaultProvider Provider

func mustGetDefaultProvider() Provider {
	if DefaultProvider == nil {
		slog.Panicf("Config read before default provider set")
	}

	return DefaultProvider
}

// Get returns the config value from the default provider at the given path
func Get(path string) Value { return mustGetDefaultProvider().Get(path) }

// Load populates the given config struct using the default provider
func Load(conf interface{}) { mustGetDefaultProvider().Load(conf) }
