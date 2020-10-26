package firehose

import "context"

// MockClient can be used by unit tests
type MockClient struct {
}

// Publish does nothing
func (m MockClient) Publish(context.Context, string, interface{}) error {
	return nil
}

// Subscribe is not implemented
func (m MockClient) Subscribe(string, Handler) {
	panic("implement me")
}
