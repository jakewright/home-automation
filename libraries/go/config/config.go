package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Provider interface {
	Has(string) bool
	Get(string) Value
}

// Config holds a nested map of config values and provides
// helper functions for easier access and type casting.
type Config struct {
	Map map[string]interface{}
}

// Value is returned from Get and has
// receiver methods for casting to various types.
type Value struct {
	raw interface{}
}

var DefaultProvider Provider
func mustGetDefaultProvider() Provider {
	if DefaultProvider == nil {
		panic("Config read before default provider set")
	}

	return DefaultProvider
}

func Has(path string) bool { return mustGetDefaultProvider().Has(path) }
func Get(path string) Value { return mustGetDefaultProvider().Get(path) }

func New(content map[string]interface{}) Provider {
	return &Config{
		Map: content,
	}
}

// Has returns whether the config has a raw at the given path e.g. "redis.host"
func (c *Config) Has(path string) bool {
	v := c.Get(path)
	return v.raw != nil
}

// Get returns the raw at the given path e.g. "redis.host"
func (c *Config) Get(path string) Value {
	return Value{
		raw: reduce(strings.Split(path, "."), c.Map),
	}
}

func reduce(parts []string, value interface{}) interface{} {
	// If this is the last part of the key
	if len(parts) == 0 {
		return value
	}

	// If raw is not a map then we can't continue
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

// Int converts the raw to an int and panics if it cannot be represented.
// The first default is returned if raw is not defined.
func (v Value) Int(defaults ...int) int {
	// Return the default if the raw is undefined
	if v.raw == nil {
		// Make sure there's at least one thing in the list
		defaults = append(defaults, 0)
		return defaults[0]
	}

	switch t := v.raw.(type) {
	case int:
		return t
	case float64:
		if t != float64(int(t)) {
			panic(fmt.Sprintf("%v cannot be represented as an int", t))
		}

		return int(t)
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			panic(err)
		}
		return i
	default:
		panic(fmt.Sprintf("%v is of unsupported type %v", t, reflect.TypeOf(t).String()))
	}
}

// String converts the raw to a string. The first default is returned if raw is not defined.
func (v Value) String(defaults ...string) string {
	if v.raw == nil {
		defaults = append(defaults, "")
		return defaults[0]
	}

	return fmt.Sprintf("%s", v.raw)
}
