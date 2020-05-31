package config

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/jakewright/home-automation/libraries/go/slog"
)

// Value is returned from Get and has
// receiver methods for casting to various types.
type Value struct {
	key    string
	raw    interface{}
	exists bool
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
			slog.Panicf("%v cannot be represented as an int", t)
		}

		return int(t)
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			slog.Panicf("failed to convert string to int: %v", err)
		}
		return i
	default:
		slog.Panicf("%v is of unsupported type %v", t, reflect.TypeOf(t).String())
	}

	return 0 // Never hit
}

// String converts the raw to a string. The first default is returned if raw is not defined.
func (v Value) String(defaults ...string) string {
	if v.raw == nil {
		defaults = append(defaults, "")
		return defaults[0]
	}

	return fmt.Sprintf("%s", v.raw)
}

// Bool converts the raw to a bool and panics if it cannot be represented.
// The first default is returned if raw is not defined.
func (v Value) Bool(defaults ...bool) bool {
	// Return the first default if the raw is undefined
	if v.raw == nil {
		// Make sure there's at least one thing in the list
		defaults = append(defaults, false)
		return defaults[0]
	}

	switch t := v.raw.(type) {
	case string:
		b, err := strconv.ParseBool(t)
		if err != nil {
			slog.Panicf("failed to parse bool: %v", err)
		}
		return b

	case bool:
		return t

	default:
		slog.Panicf("%v is of unsupported type %v", t, reflect.TypeOf(t).String())
	}

	return false
}

// Duration converts the raw to a time.Duration.
// The first default is returned if raw is not defined.
func (v Value) Duration(defaults ...time.Duration) time.Duration {
	// Return the first default if raw is undefined
	if v.raw == nil {
		// Make sure there's at least one thing in the list
		defaults = append(defaults, 0)
		return defaults[0]
	}

	switch t := v.raw.(type) {
	case string:
		d, err := time.ParseDuration(t)
		if err != nil {
			slog.Panicf("Failed to parse duration: %v", err)
		}
		return d
	default:
		slog.Panicf("%v is of unsupported type %v", t, reflect.TypeOf(t).String())
	}

	return 0
}
