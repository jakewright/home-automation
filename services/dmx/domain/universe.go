package domain

// UniverseNumber is the OLA universe number that can range from 1 to 65535
// https://www.openlighting.org/ola/get-help/ola-faq/#Universes
type UniverseNumber uint16

// Universe represents a 512 channel space
type Universe struct {
	// values holds the current value of all channels. This should not be read
	// directly as it may be out of date. It is updated with the latest values
	// from the fixtures when DMXValues() is called.
	values [512]byte

	// fixtures is a set of fixtures in the universe. Note that this does not
	// need to be a complete set. Only the fixtures that you care about
	// changing need to be added to a universe. Values of all other channels
	// will be maintained.
	fixtures []Fixture
}

// NewUniverse returns a new universe containing the given fixtures. Each
// fixture is hydrated with the relevant slice of values from the byte array.
func NewUniverse(values [512]byte, fixtures ...Fixture) (*Universe, error) {
	// Hydrate each fixture
	for _, f := range fixtures {
		slice := values[f.offset() : f.offset()+f.length()]
		if err := f.hydrate(slice); err != nil {
			return nil, err
		}
	}

	return &Universe{
		values:   values,
		fixtures: fixtures,
	}, nil
}

// DMXValues returns the value of all channels in the universe
func (u *Universe) DMXValues() [512]byte {
	for _, f := range u.fixtures {
		copy(u.values[f.offset():], f.dmxValues())
	}
	return u.values
}
