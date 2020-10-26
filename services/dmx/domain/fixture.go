package domain

import (
	"sort"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/oops"
)

// Fixture types
const (
	FixtureTypeMegaParProfile = "mega_par_profile"
)

// Fixture is the interface that all DMX
// fixture structs must implement
type Fixture interface {
	// ID returns the device ID of the fixture
	ID() string

	// UniverseNumber returns the number of the
	// universe of which the fixture is a part
	UniverseNumber() UniverseNumber

	// offset returns the device's offset
	// into the universe's channel space
	offset() int

	// length returns the number of channels
	// that the fixture occupies
	length() int

	// hydrate takes a slice of DMX values of
	// size Length() and sets internal state
	hydrate([]byte) error

	// dmxValues returns a slice of size Length() of values
	// representing the current state of the fixture
	dmxValues() []byte

	// setHeader is used by newFromDeviceHeader to set
	// properties common to all fixtures
	setHeader(header *devicedef.Header) error
}

// NewFixture returns a Fixture based on the device's fixture type attribute
func NewFixture(h *devicedef.Header) (Fixture, error) {
	fixtureType, ok := h.Attributes["fixture_type"].(string)
	if !ok {
		return nil, oops.PreconditionFailed("fixture_type not found in %s device header", h.Id)
	}

	var f Fixture

	switch fixtureType {
	case FixtureTypeMegaParProfile:
		f = &MegaParProfile{}
	default:
		return nil, oops.InternalService("device %s has invalid fixture type '%s'", h.Id, fixtureType)
	}

	if err := f.setHeader(h); err != nil {
		return nil, err
	}

	return f, nil
}

// ValidateFixtures checks for fixtures with overlapping channel ranges
func ValidateFixtures(fs []Fixture) error {
	// Sort the fixtures by offset
	sort.Slice(fs, func(i, j int) bool {
		return fs[i].offset() < fs[j].offset()
	})

	// Make sure each fixture ends before the next one begins
	for i := 0; i < len(fs)-1; i++ {
		if fs[i].offset()+fs[i].length() > fs[i+1].offset() {
			return oops.InternalService("universe has overlapping fixtures")
		}
	}

	return nil
}
