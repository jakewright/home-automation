package routes

import (
	"context"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/services/device-registry/def"
)

// ListDevices lists all devices known by the registry. Results can be filtered by controller name.
func (c *Controller) ListDevices(ctx context.Context, body *deviceregistrydef.ListDevicesRequest) (*deviceregistrydef.ListDevicesResponse, error) {
	var devices []*devicedef.Header
	var err error

	if controllerName, set := body.GetControllerName(); set {
		devices, err = c.DeviceRepository.FindByController(controllerName)
	} else {
		devices, err = c.DeviceRepository.FindAll()
	}
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find devices")
	}

	return &deviceregistrydef.ListDevicesResponse{
		DeviceHeaders: devices,
	}, nil
}

// GetDevice returns a specific device by ID
func (c *Controller) GetDevice(ctx context.Context, body *deviceregistrydef.GetDeviceRequest) (*deviceregistrydef.GetDeviceResponse, error) {
	device, err := c.DeviceRepository.Find(body.GetDeviceId())
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find device %q", body.GetDeviceId())
	}
	if device == nil {
		return nil, oops.NotFound("device %q not found", body.GetDeviceId())
	}

	return &deviceregistrydef.GetDeviceResponse{
		DeviceHeader: device,
	}, nil
}
