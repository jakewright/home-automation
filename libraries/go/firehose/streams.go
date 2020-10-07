package firehose

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v7"
)

// StreamsEvent is a message received from a Redis stream
type StreamsEvent struct {
	channel string
	payload []byte
}

// Channel returns the name of the stream from which the
// message was received
func (e *StreamsEvent) Channel() string {
	return e.channel
}

// Decode unmarshals the raw payload into v
func (e *StreamsEvent) Decode(v interface{}) error {
	return json.Unmarshal(e.payload, v)
}

// StreamsClient uses Redis Streams to publish and subscribe
type StreamsClient struct {
	client *redis.Client
}

// NewStreamsClient returns a new StreamsClient that uses
// the given Redis client to interact with the streams API
func NewStreamsClient(client *redis.Client) *StreamsClient {
	return &StreamsClient{
		client: client,
	}
}

// Publish adds the message to the channel (stream) with the
// given name. The message is marshaled to JSON and added
// under a key called "data".
func (s *StreamsClient) Publish(channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.client.XAdd(&redis.XAddArgs{
		Stream:       channel,
		MaxLenApprox: 100,
		ID:           "*", // Generate an ID
		Values: map[string]interface{}{
			"data": data,
		},
	}).Err()
}

// Subscribe registers a handler for the given channel
func (s *StreamsClient) Subscribe(channel string, handler Handler) {
	panic("implement me")
}

// GetName returns the friendly name for the process
func (s *StreamsClient) GetName() string {
	return "Firehose"
}

// Start subscribes to the channels and listens for messages
func (s *StreamsClient) Start(ctx context.Context) error {
	return nil
}
