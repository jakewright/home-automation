package util

import "github.com/jakewright/home-automation/libraries/go/config"

// IsProd returns whether the current environment is production, based on config.
func IsProd() bool {
	return config.Get("environment").String("prod") == "prod"
}
