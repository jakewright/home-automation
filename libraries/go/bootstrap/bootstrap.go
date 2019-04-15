package bootstrap

import (
	"fmt"
	"home-automation/libraries/go/config"
	"home-automation/libraries/go/firehose"
	"home-automation/libraries/go/http"
	"home-automation/libraries/go/slog"
	"os"

	"github.com/go-redis/redis"
)

// Boot performs standard service startup tasks
func Init(serviceName string) error {
	// Create default HTTP client
	apiGateway := os.Getenv("API_GATEWAY")
	if apiGateway == "" {
		return fmt.Errorf("API_GATEWAY env var not set")
	}
	httpClient, err := http.New(apiGateway, "data")
	if err != nil {
		return err
	}
	http.DefaultClient = httpClient

	// Load config
	var configRsp map[string]interface{}
	_, err = http.Get(fmt.Sprintf("service.config/read/%s", serviceName), &configRsp)
	if err != nil {
		return err
	}
	config.DefaultProvider = config.New(configRsp)

	// Connect to Redis
	if config.Has("redis.host") {
		host := config.Get("redis.host").String()
		port := config.Get("redis.port").Int()
		addr := fmt.Sprintf("%s:%d", host, port)
		slog.Info("Connecting to Redis at address %s", addr)
		firehose.DefaultPublisher = firehose.New(redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}))
	}

	return nil
}
