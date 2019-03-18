package firehose

import "github.com/go-redis/redis"

type Publisher interface {
	Publish(string, string) error
}

var DefaultPublisher Publisher
func mustGetDefaultPublisher() Publisher {
	if DefaultPublisher == nil {
		panic("Firehose used before default publisher set")
	}

	return DefaultPublisher
}

func Publish(channel, message string) error { return mustGetDefaultPublisher().Publish(channel, message) }

type RedisClient struct {
	client *redis.Client
}

func New(client *redis.Client) Publisher {
	return &RedisClient{
		client: client,
	}
}

func (c *RedisClient) Publish(channel, message string) error {
	return c.client.Publish(channel, message).Err()
}