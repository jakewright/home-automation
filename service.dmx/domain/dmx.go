package domain

import (
	"sort"

	"github.com/jakewright/home-automation/libraries/go/errors"
	deviceregistryproto "github.com/jakewright/home-automation/service.device-registry/proto"
)

// DeviceHeader is a wrapper that adds typed Attributes
type DeviceHeader struct {
	*deviceregistryproto.DeviceHeader
	Attributes Attributes `json:"attributes"`
}

// Attributes describe a fixture
type Attributes struct {
	FixtureType string `json:"fixture_type"`
	Offset      int    `json:"offset"`
}

// Universe represents a set of fixtures in a 512 channel space
type Universe struct {
	Number   int
	Fixtures []Fixture
}

// Valid returns false if any fixtures have overlapping channel ranges
func (u *Universe) Valid() bool {
	// Don't modify the slice
	var f []Fixture
	copy(f, u.Fixtures)

	// Sort the fixtures by offset
	sort.Slice(f, func(i, j int) bool {
		return f[i].Offset() < f[j].Offset()
	})

	// Make sure each fixture ends before the next one begins
	for i := 0; i < len(f)-1; i++ {
		if f[i].Offset()+len(f[i].DMXValues()) > f[i+1].Offset() {
			return false
		}
	}

	return true
}

// DMXValues returns the value of all channels in the universe
func (u *Universe) DMXValues() [512]byte {
	var v [512]byte
	for _, f := range u.Fixtures {
		copy(v[f.Offset():], f.DMXValues())
	}
	return v
}

// Fixture is an addressable device
type Fixture interface {
	ID() string
	DMXValues() []byte
	Offset() int
	SetProperty(string, string) error
}

// NewFixtureFromDeviceHeader returns a Fixture based on the device's fixture type attribute
func NewFixtureFromDeviceHeader(h *DeviceHeader) (Fixture, error) {
	switch h.Attributes.FixtureType {
	case "mega_par_profile":
		return &MegaParProfile{DeviceHeader: h}, nil
	}
	return nil, errors.InternalService("Device %s has invalid fixture type '%s'", h.ID, h.Attributes.FixtureType)
}
