package bootstrap

import (
	"context"
	"fmt"
	"home-automation/libraries/go/api"
	"home-automation/libraries/go/config"
	"home-automation/libraries/go/firehose"
	"home-automation/libraries/go/router"
	"home-automation/libraries/go/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
)

// Boot performs standard service startup tasks
func Init(serviceName string) error {
	// Create default API client
	apiGateway := os.Getenv("API_GATEWAY")
	if apiGateway == "" {
		return fmt.Errorf("API_GATEWAY env var not set")
	}
	apiClient, err := api.New(apiGateway, "data")
	if err != nil {
		return err
	}
	api.DefaultClient = apiClient

	// Load config
	var configRsp map[string]interface{}
	_, err = api.Get(fmt.Sprintf("service.config/read/%s", serviceName), &configRsp)
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

func Run() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	errs := make(chan error)
	go func() {
		errs <- router.ListenAndServe()
	}()

	select {
	case err := <-errs:
		slog.Error("Router unexpectedly stopped: %v", err)
		os.Exit(1)
	case s := <-sig:
		slog.Info("Received signal %v; exiting...", s)
	}

	// A short timeout because Docker will kill us after 10 seconds anyway
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown gracefully: %v", err)
	}

	slog.Info("Service stopped")
}
