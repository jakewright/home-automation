package firehose

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
