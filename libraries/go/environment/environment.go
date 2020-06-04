package environment

import "github.com/jakewright/home-automation/libraries/go/config"

const envProd = "prod"

var conf struct {
	Environment string `envconfig:"default=prod"`
}

func init() {
	config.Load(&conf)
}

// IsProd returns whether the current environment is production, based on config.
func IsProd() bool {
	return conf.Environment == envProd
}
