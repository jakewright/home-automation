package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/healthz"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

// getRedisClient returns a cached instance of a redis.Client.
// If it is being called for the first time, a new connection to
// Redis is initiated. Connection options are read from config.
func (s *Service) getRedisClient() (*redis.Client, error) {
	if s.redisClient == nil {
		conf := struct {
			RedisHost string
			RedisPort int
		}{}
		config.Load(&conf)

		addr := fmt.Sprintf("%s:%d", conf.RedisHost, conf.RedisPort)
		slog.Infof("Connecting to Redis at address %s", addr)
		s.redisClient = redis.NewClient(&redis.Options{
			Addr:            addr,
			Password:        "",
			DB:              0,
			MaxRetries:      5,
			MinRetryBackoff: time.Second,
			MaxRetryBackoff: time.Second * 5,
		})

		s.runner.addDeferred(func() error {
			err := s.redisClient.Close()
			if err != nil {
				slog.Errorf("Failed to close Redis connection: %v", err)
			} else {
				slog.Debugf("Closed Redis connection")
			}
			return err
		})

		healthCheck := func(ctx context.Context) error {
			_, err := s.redisClient.Ping(ctx).Result()
			return err
		}

		healthz.RegisterCheck("redis", healthCheck)

		if err := healthCheck(context.Background()); err != nil {
			return nil, err
		}
	}

	return s.redisClient, nil
}
