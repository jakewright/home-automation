package universe

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
	"github.com/jakewright/home-automation/service.dmx/domain"
)

// Loader loads device metadata and instantiates fixtures
type Loader struct {
	ServiceName string
	Universe    *Universe
}

// FetchDevices loads devices from the device registry, creates fixtures, and adds them to the universe.
func (l *Loader) FetchDevices(ctx context.Context) error {
	rsp, err := (&deviceregistrydef.ListDevicesRequest{
		ControllerName: l.ServiceName,
	}).Do(ctx)
	if err != nil {
		return oops.WithMessage(err, "failed to fetch devices")
	}

	for _, device := range rsp.DeviceHeaders {
		switch {
		case device.ControllerName != l.ServiceName:
			return oops.InternalService("device %s is not for this controller", device.Id)
		case device.Type != "dmx":
			return oops.InternalService("device %s does not have type dmx", device.Id)
		}

		fixture, err := domain.NewFixtureFromDeviceHeader(device)
		if err != nil {
			return oops.WithMessage(err, "failed to create fixture")
		}

		l.Universe.AddFixture(fixture)
	}

	return nil
}
