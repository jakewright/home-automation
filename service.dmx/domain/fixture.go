package domain

import (
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

// Fixture types
const (
	FixtureTypeMegaParProfile = "mega_par_profile"
)

// Fixture is an addressable device
type Fixture interface {
	SetHeader(*deviceregistrydef.DeviceHeader) error
	ID() string
	ToDef() *devicedef.Device
	DMXValues() []byte
	Offset() int
	SetProperties(map[string]interface{}) (bool, error)
	Copy() Fixture
}

// NewFixtureFromDeviceHeader returns a Fixture based on the device's fixture type attribute
func NewFixtureFromDeviceHeader(h *deviceregistrydef.DeviceHeader) (Fixture, error) {
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

	if err := f.SetHeader(h); err != nil {
		return nil, err
	}

	return f, nil
}
