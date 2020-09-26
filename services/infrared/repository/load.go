package repository

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/services/device-registry/def"
	"github.com/jakewright/home-automation/services/infrared/domain"
)

// Loader loads device metadata and instantiates devices
type Loader struct {
	ServiceName string
	Repository  *DeviceRepository
}

func (l *Loader) FetchDevices(ctx context.Context) error {
	rsp, err := (&deviceregistrydef.ListDevicesRequest{
		ControllerName: l.ServiceName,
	}).Do(ctx)
	if err != nil {
		return oops.WithMessage(err, "failed to fetch devices")
	}

	for _, device := range rsp.DeviceHeaders {
		// Sanity check the controller name
		if device.ControllerName != l.ServiceName {
			return oops.InternalService("device %s is not for this controller", device.Id)
		}

		device, err := domain.NewDeviceFromDeviceHeader(device)
		if err != nil {
			return oops.WithMessage(err, "failed to create device")
		}

		l.Repository.AddDevice(device)
	}

	return nil
}
