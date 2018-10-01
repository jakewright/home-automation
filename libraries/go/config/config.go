package config

import (
	"strings"
)

type Config struct {
	Map map[string]interface{}
}

type configValue struct {
	value interface{}
}

func (c *Config) Has(path string) bool {
	v := c.Get(path)
	return v.value != nil
}

func (c *Config) Get(path string) configValue {
	return configValue{
		value: reduce(strings.Split(path, "."), c.Map),
	}
}

func reduce(parts []string, value interface{}) interface{} {
	// If this is the last part of the key
	if len(parts) == 0 {
		return value
	}

	// If value is not a map then we can't continue
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}

	// If the key we are searching for is not defined
	value, ok = valueMap[parts[0]]
	if !ok {
		return nil
	}

	return reduce(parts[1:], value)
}

func (v configValue) Int(defaults ...int) int {
	if len(defaults) == 0 {
		defaults = append(defaults, 0)
	}

	switch t := v.value.(type) {
	case int:
		return t
	case float64:
		return int(t)
	default:
		return defaults[0]
	}
}

func (v configValue) String(defaults ...string) string {
	if len(defaults) == 0 {
		defaults = append(defaults, "")
	}

	s, ok := v.value.(string)
	if !ok {
		return defaults[0]
	}
	return s
}
