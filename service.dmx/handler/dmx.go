package handler

import (
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/errors"
	dmxdef "github.com/jakewright/home-automation/service.dmx/def"
	"github.com/jakewright/home-automation/service.dmx/universe"
)

type setter interface {
	Set(universe int, values [512]byte) error
}

// DMXController handles device requests
type DMXController struct {
	Universe *universe.Universe
	Setter   setter
}

// Read returns the current state of a fixture
func (c *DMXController) Read(r *Request, body *dmxdef.GetDeviceRequest) (*dmxdef.GetDeviceResponse, error) {
	fixture := c.Universe.Find(body.DeviceId)
	if fixture == nil {
		return nil, errors.NotFound("device %q not found", body.DeviceId)
	}

	return &dmxdef.GetDeviceResponse{
		Device: fixture.ToDef(),
	}, nil
}

// Update modifies fixture properties
func (c *DMXController) Update(r *Request, body *dmxdef.UpdateDeviceRequest) (*dmxdef.UpdateDeviceResponse, error) {
	errParams := map[string]string{
		"device_id": body.DeviceId,
	}

	fixture := c.Universe.Find(body.DeviceId)
	if fixture == nil {
		return nil, errors.NotFound("device %q not found", body.DeviceId, errParams)
	}

	changed, err := fixture.SetProperties(body.State)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update fixture", errParams)
	}

	if err := c.Setter.Set(c.Universe.Number, c.Universe.DMXValues()); err != nil {
		return nil, errors.WithMessage(err, "failed to set DMX values", errParams)
	}

	if changed {
		if err := (&devicedef.DeviceStateChangedEvent{
			Device: fixture.ToDef(),
		}).Publish(); err != nil {
			return nil, errors.WithMessage(err, "failed to publish state changed event", errParams)
		}
	}

	return &dmxdef.UpdateDeviceResponse{
		Device: fixture.ToDef(),
	}, nil
}
