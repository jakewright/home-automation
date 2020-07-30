package dmx

import "context"

// MockGetSetter can be used in tests
type MockGetSetter struct {
	Values [512]byte
}

var _ GetSetter = (*MockGetSetter)(nil)

// GetValues returns the current values
func (m *MockGetSetter) GetValues(ctx context.Context) ([512]byte, error) {
	return m.Values, nil
}

// SetValues replaces the values with a new slice
func (m *MockGetSetter) SetValues(ctx context.Context, values [512]byte) error {
	m.Values = values
	return nil
}
