package domain

import (
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/errors"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

// Fixture types
const (
	FixtureTypeMegaParProfile = "mega_par_profile"
)

// Fixture is an addressable device
type Fixture interface {
	ID() string
	ToDef() *devicedef.Device
	DMXValues() []byte
	Offset() int
	SetProperties(map[string]interface{}) (bool, error)
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
	case FixtureTypeMegaParProfile:
		return &MegaParProfile{DeviceHeader: h}, nil
	}
	return nil, errors.InternalService("device %s has invalid fixture type '%s'", h.Id, h.Attributes.FixtureType)
}
