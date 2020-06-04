package handler

import (
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/dsync"
	"github.com/jakewright/home-automation/libraries/go/oops"
	dmxproxydef "github.com/jakewright/home-automation/service.dmx-proxy/def"
	dmxdef "github.com/jakewright/home-automation/service.dmx/def"
	"github.com/jakewright/home-automation/service.dmx/universe"
)

// Controller handles requests
type Controller struct {
	Universe *universe.Universe
}

// GetDevice returns the current state of a fixture
func (c *Controller) GetDevice(r *request, body *dmxdef.GetDeviceRequest) (*dmxdef.GetDeviceResponse, error) {
	fixture := c.Universe.Find(body.DeviceId)
	if fixture == nil {
		return nil, oops.NotFound("device %q not found", body.DeviceId)
	}

	return &dmxdef.GetDeviceResponse{
		Device: fixture.ToDef(),
	}, nil
}

// UpdateDevice modifies fixture properties
func (c *Controller) UpdateDevice(r *request, body *dmxdef.UpdateDeviceRequest) (*dmxdef.UpdateDeviceResponse, error) {
	errParams := map[string]string{
		"device_id": body.DeviceId,
	}

	// Take a lock on the entire universe because even though we're only
	// updating a single device, we need to send the entire universe's state
	// over the wire to the fixtures. We therefore don't want simultaneous
	// update requests to interleave.
	lock, err := dsync.Lock(r, "dmx-universe", c.Universe.Number)
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}
	defer lock.Unlock()

	fixture := c.Universe.Find(body.DeviceId)
	if fixture == nil {
		return nil, oops.NotFound("device %q not found", body.DeviceId, errParams)
	}

	changed, err := fixture.SetProperties(body.State)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to update fixture", errParams)
	}

	values := c.Universe.DMXValues(fixture)
	if _, err = (&dmxproxydef.SetRequest{
		Universe: c.Universe.Number,
		Values:   values[:],
	}).Do(r); err != nil {
		return nil, oops.WithMessage(err, "failed to set DMX values", errParams)
	}

	if changed {
		if err := (&devicedef.DeviceStateChangedEvent{
			Device: fixture.ToDef(),
		}).Publish(); err != nil {
			return nil, oops.WithMessage(err, "failed to publish state changed event", errParams)
		}
	}

	c.Universe.Save(fixture)

	return &dmxdef.UpdateDeviceResponse{
		Device: fixture.ToDef(),
	}, nil
}
