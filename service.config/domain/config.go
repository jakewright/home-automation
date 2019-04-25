package domain

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/jakewright/home-automation/libraries/go/errors"

	"gopkg.in/yaml.v2"
)

// Config is an abstraction around the map that holds config values
type Config struct {
	config map[string]interface{}
	lock   sync.RWMutex
}

// SetFromBytes sets the internal config based on a byte array of YAML
func (c *Config) SetFromBytes(data []byte) (bool, error) {
	var rawConfig interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return false, err
	}

	configUntyped, ok := rawConfig.(map[interface{}]interface{})
	if !ok {
		return false, fmt.Errorf("config is not a map")
	}

	// YAML allows non-string map keys but we need to be able to marshal the config into JSON which only allows strings
	config, err := convertKeysToString(configUntyped)
	if err != nil {
		return false, err
	}

	if reflect.DeepEqual(config, c.config) {
		return false, nil
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.config = config

	return true, nil
}

// Get returns the config for the given service
func (c *Config) Get(serviceName string) (map[string]interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	errParams := map[string]string{"serviceName": serviceName}

	a, ok := c.config["base"].(map[string]interface{})
	if !ok {
		return nil, errors.InternalService("base config is not a map", errParams)
	}

	// If no config is defined for the service
	if _, ok = c.config[serviceName]; !ok {
		// Return the base config
		return a, nil
	}

	b, ok := c.config[serviceName].(map[string]interface{})
	if !ok {
		return nil, errors.InternalService("service %q config is not a map", errParams)
	}

	// Merge the maps with service config taking precedence
	config := make(map[string]interface{})
	for k, v := range a {
		config[k] = v
	}
	for k, v := range b {
		config[k] = v
	}

	return config, nil
}

// convertKeysToString recursively iterates over a map with interface{} keys and asserts that they are strings
func convertKeysToString(m map[interface{}]interface{}) (map[string]interface{}, error) {
	n := make(map[string]interface{})

	for k, v := range m {
		// Assert that the key is a string
		str, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("config key is not a string")
		}

		if vMap, ok := v.(map[interface{}]interface{}); ok {
			var err error
			v, err = convertKeysToString(vMap)
			if err != nil {
				return nil, err
			}
		}

		n[str] = v
	}

	return n, nil
}
