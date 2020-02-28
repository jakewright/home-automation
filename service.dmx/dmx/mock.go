package dmx

// Mock can be used in unit tests
type Mock struct {
	Universe int
	Values   [512]byte
}

// Set sets all of the DMX values for the given universe
func (m *Mock) Set(universe int, values [512]byte) error {
	m.Universe = universe
	m.Values = values

	return nil
}
