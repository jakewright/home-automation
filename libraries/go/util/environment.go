package util

import "github.com/jakewright/home-automation/libraries/go/config"

const envProd = "prod"

var env = config.Get("ENV").String(envProd)

// IsProd returns whether the current environment is production, based on config.
func IsProd() bool {
	return env == envProd
}
