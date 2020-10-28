package routes

import (
	"context"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/distsync"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/libraries/go/oops"
	def "github.com/jakewright/home-automation/services/dmx/def"
	dmxdef "github.com/jakewright/home-automation/services/dmx/def"
	"github.com/jakewright/home-automation/services/dmx/dmx"
	"github.com/jakewright/home-automation/services/dmx/domain"
	"github.com/jakewright/home-automation/services/dmx/repository"
)

// Controller handles requests
type Controller struct {
	Repository *repository.FixtureRepository
	Client     *dmx.Client
	Publisher  firehose.Publisher
}

// GetMegaParProfile returns the current state of a device of type mega-par-profile
func (c *Controller) GetMegaParProfile(ctx context.Context, body *dmxdef.GetMegaParProfileRequest) (*def.MegaParProfileResponse, error) {
	errParams := map[string]string{
		"device_id": body.GetDeviceId(),
	}

	lock, err := distsync.Lock(ctx, "device", body.GetDeviceId())
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}
	defer lock.Unlock()

	f := c.Repository.Find(body.GetDeviceId())
	if f == nil {
		return nil, oops.NotFound("device %q not found", body.GetDeviceId())
	}

	megaParProfile, ok := f.(*domain.MegaParProfile)
	if !ok {
		return nil, oops.BadRequest("device %q is not a MegaParProfile", body.DeviceId)
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

	return &dmxdef.MegaParProfileResponse{
		Header: megaParProfile.Header,
		Properties: map[string]*devicedef.Property{
			"power":      (&devicedef.Property{}).SetType("bool"),
			"brightness": (&devicedef.Property{}).SetType("uint8"),
			"color":      (&devicedef.Property{}).SetType("rgb"),
			"strobe":     (&devicedef.Property{}).SetType("uint8"),
		},
		State: megaParProfile.State(),
	}, nil
}

// UpdateMegaParProfile updates a device of type mega-par-profile
func (c *Controller) UpdateMegaParProfile(
	ctx context.Context,
	body *dmxdef.UpdateMegaParProfileRequest,
) (*dmxdef.MegaParProfileResponse, error) {
	errParams := map[string]string{
		"device_id": body.GetDeviceId(),
	}

	lock, err := distsync.Lock(ctx, "device", body.GetDeviceId())
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}
	defer lock.Unlock()

	f := c.Repository.Find(body.GetDeviceId())
	if f == nil {
		return nil, oops.NotFound("device %q not found", body.GetDeviceId())
	}

	megaParProfile, ok := f.(*domain.MegaParProfile)
	if !ok {
		return nil, oops.BadRequest("device %q is not a MegaParProfile", body.DeviceId)
	}

	values, err := c.Client.GetValues(ctx, f.UniverseNumber())
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}

	u, err := domain.NewUniverse(values, f)
	if err != nil {
		return nil, oops.WithMetadata(err, errParams)
	}

	megaParProfile.ApplyState(body.State)

	values = u.DMXValues()
	if err = c.Client.SetValues(ctx, f.UniverseNumber(), values); err != nil {
		return nil, oops.WithMessage(err, "failed to set DMX values", errParams)
	}

	// TODO: enable event publishing
	// if err := (&devicedef.DeviceStateChangedEvent{
	// 	Header: megaParProfile.Header,
	// 	State:  megaParProfile.State(),
	// }).Publish(ctx, c.Publisher); err != nil {
	// 	return nil, oops.WithMessage(err, "failed to publish state changed event", errParams)
	// }

	return &dmxdef.MegaParProfileResponse{
		Header: megaParProfile.Header,
		Properties: map[string]*devicedef.Property{
			"power":      (&devicedef.Property{}).SetType("bool"),
			"brightness": (&devicedef.Property{}).SetType("uint8"),
			"color":      (&devicedef.Property{}).SetType("rgb"),
			"strobe":     (&devicedef.Property{}).SetType("uint8"),
		},
		State: megaParProfile.State(),
	}, nil
}
