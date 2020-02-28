package firehose

// MockClient can be used by unit tests
type MockClient struct {
}

// Publish does nothing
func (m MockClient) Publish(string, interface{}) error {
	return nil
}

// Subscribe is not implemented
func (m MockClient) Subscribe(string, RawHandlerFunc) {
	panic("implement me")
}

// config is not implemented
func (m MockClient) config() *Config {
	panic("implement me")
}
