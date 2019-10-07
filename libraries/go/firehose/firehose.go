package firehose

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

// Publisher is the interface that wraps the Publish function
type Publisher interface {
	Publish(string, interface{}) error
}

// DefaultPublisher is a global instance of Publisher
var DefaultPublisher Publisher

func mustGetDefaultPublisher() Publisher {
	if DefaultPublisher == nil {
		panic("Firehose used before default publisher set")
	}

	return DefaultPublisher
}

// Publish sends the given message on the given channel using the default publisher
func Publish(channel string, message interface{}) error {
	return mustGetDefaultPublisher().Publish(channel, message)
}

// RedisClient wraps a redis.Client and exposes a Publish method
type RedisClient struct {
	client *redis.Client
}

// New returns a RedisClient
func New(client *redis.Client) Publisher {
	return &RedisClient{
		client: client,
	}
}

// Publish emits the given message on the given redis channel
func (c *RedisClient) Publish(channel string, message interface{}) error {
	// The go-redis Publish() function will error if the message is not a string,
	// []byte, bool, or numeric type. Strings will be converted to []byte.

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.client.Publish(channel, data).Err()
}
