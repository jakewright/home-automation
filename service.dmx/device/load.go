package device

import (
	"context"
	"fmt"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/rpc"
	"github.com/jakewright/home-automation/service.dmx/domain"
)

// Loader loads device metadata and instantiates fixtures
type Loader struct {
	ServiceName string
	Universe    *domain.Universe
}

// FetchDevices loads devices from the device registry, creates fixtures, and adds them to the universe.
func (l *Loader) FetchDevices(ctx context.Context) error {
	url := fmt.Sprintf("service.device-registry/devices?controller_name=%s", l.ServiceName)
	var rsp []*domain.DeviceHeader
	if _, err := rpc.Get(ctx, url, &rsp); err != nil {
		return err
	}

	for _, device := range rsp {
		switch {
		case device.ControllerName != l.ServiceName:
			return errors.InternalService("device %s is not for this controller", device.Id)
		case device.Type != "dmx":
			return errors.InternalService("device %s does not have type dmx", device.Id)
		}

		fixture, err := domain.NewFixtureFromDeviceHeader(device)
		if err != nil {
			return errors.WithMessage(err, "failed to create fixture")
		}

		l.Universe.AddFixture(fixture)
	}

	return nil
}
