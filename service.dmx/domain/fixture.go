package domain

import (
	"github.com/jakewright/home-automation/libraries/go/errors"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

// Fixture is an addressable device
type Fixture interface {
	ID() string
	DMXValues() []byte
	Offset() int
	SetProperties([]byte) (bool, error)
}

// DeviceHeader is a wrapper that adds typed Attributes
type DeviceHeader struct {
	*deviceregistrydef.DeviceHeader
	Attributes Attributes `json:"attributes"`
}

// Attributes describe a fixture
type Attributes struct {
	FixtureType string `json:"fixture_type"`
	Offset      int    `json:"offset"`
}

// NewFixtureFromDeviceHeader returns a Fixture based on the device's fixture type attribute
func NewFixtureFromDeviceHeader(h *DeviceHeader) (Fixture, error) {
	switch h.Attributes.FixtureType {
	case "mega_par_profile":
		return &MegaParProfile{DeviceHeader: h}, nil
	}
	return nil, errors.InternalService("Device %s has invalid fixture type '%s'", h.Id, h.Attributes.FixtureType)
}
