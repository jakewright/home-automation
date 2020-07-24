package handler

import (
	"context"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/dsync"
	"github.com/jakewright/home-automation/libraries/go/oops"
	dmxdef "github.com/jakewright/home-automation/service.dmx/def"
	"github.com/jakewright/home-automation/service.dmx/dmx"
	"github.com/jakewright/home-automation/service.dmx/domain"
	"github.com/jakewright/home-automation/service.dmx/repository"
)

// Controller handles requests
type Controller struct {
	Repository *repository.FixtureRepository
	Client     *dmx.Client
}

// GetDevice returns the current state of a fixture
func (c *Controller) GetDevice(ctx context.Context, body *dmxdef.GetDeviceRequest) (*dmxdef.GetDeviceResponse, error) {
	errParams := map[string]string{
		"device_id": body.DeviceId,
	}

	lock, err := dsync.Lock(ctx, "device", body.DeviceId)
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}
	defer lock.Unlock()

	f := c.Repository.Find(body.DeviceId)
	if f == nil {
		return nil, oops.NotFound("device %q not found", body.DeviceId)
	}

	values, err := c.Client.GetValues(ctx, f.UniverseNumber())
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}

	// Instantiating a universe will hydrate the fixture
	_, err = domain.NewUniverse(values, f)
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}

	return &dmxdef.GetDeviceResponse{
		Device: f.ToDevice(),
	}, nil
}

// UpdateDevice modifies fixture properties
func (c *Controller) UpdateDevice(ctx context.Context, body *dmxdef.UpdateDeviceRequest) (*dmxdef.UpdateDeviceResponse, error) {
	errParams := map[string]string{
		"device_id": body.DeviceId,
	}

	lock, err := dsync.Lock(ctx, "device", body.DeviceId)
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}
	defer lock.Unlock()

	f := c.Repository.Find(body.DeviceId)
	if f == nil {
		return nil, oops.NotFound("device %q not found", body.DeviceId)
	}

	values, err := c.Client.GetValues(ctx, f.UniverseNumber())
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}

	u, err := domain.NewUniverse(values, f)
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}

	if err := f.SetProperties(body.State); err != nil {
		return nil, oops.WithMessage(err, "failed to update fixture", errParams)
	}

	values = u.DMXValues()
	if err = c.Client.SetValues(ctx, f.UniverseNumber(), values); err != nil {
		return nil, oops.WithMessage(err, "failed to set DMX values", errParams)
	}

	if err := (&devicedef.DeviceStateChangedEvent{
		Device: f.ToDevice(),
	}).Publish(); err != nil {
		return nil, oops.WithMessage(err, "failed to publish state changed event", errParams)
	}

	return &dmxdef.UpdateDeviceResponse{
		Device: f.ToDevice(),
	}, nil
}
