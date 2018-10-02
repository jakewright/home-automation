package bootstrap

import (
	"fmt"
	"github.com/go-redis/redis"
	"home-automation/libraries/go/client"
	"home-automation/libraries/go/config"
	"log"
	"os"
)

// Service gives access to common features
type Service struct {
	ControllerName string
	Config         *config.Config
	Redis          *redis.Client
}

// Boot performs standard Service startup tasks
func Boot(controllerName string) (*Service, error) {
	service := &Service{
		ControllerName: controllerName,
	}

	// Create API Client
	apiGateway := os.Getenv("API_GATEWAY")
	if apiGateway == "" {
		return nil, fmt.Errorf("API_GATEWAY env var not set")
	}
	apiClient, err := client.New(apiGateway)
	if err != nil {
		return nil, err
	}

	// Load config
	var configRsp map[string]interface{}
	_, err = apiClient.Get(fmt.Sprintf("Service.config/read/%s", controllerName), &configRsp)
	if err != nil {
		return nil, err
	}
	service.Config = &config.Config{
		Map: configRsp,
	}

	// Connect to Redis
	if service.Config.Has("redis.host") {
		host := service.Config.Get("redis.host").String()
		port := service.Config.Get("redis.port").Int()
		addr := fmt.Sprintf("%s:%d", host, port)
		log.Printf("Connecting to Redis at address %s\n", addr)
		service.Redis = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		})
	}

	return service, nil
}
