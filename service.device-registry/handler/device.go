package handler

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

// ListDevices lists all devices known by the registry. Results can be filtered by controller name.
func (c *Controller) ListDevices(ctx context.Context, body *deviceregistrydef.ListDevicesRequest) (*deviceregistrydef.ListDevicesResponse, error) {
	var devices []*deviceregistrydef.DeviceHeader
	var err error
	if body.ControllerName != "" {
		devices, err = c.DeviceRepository.FindByController(body.ControllerName)
	} else {
		devices, err = c.DeviceRepository.FindAll()
	}
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find devices")
	}

	// Decorate the devices with their rooms
	for _, device := range devices {
		room, err := c.RoomRepository.Find(device.RoomId)
		if err != nil {
			return nil, oops.WithMessage(err, "failed to find room %q", device.RoomId)
		}
		if room == nil {
			return nil, oops.NotFound("room %q not found", device.RoomId)
		}

		device.Room = room
	}

	return &deviceregistrydef.ListDevicesResponse{
		DeviceHeaders: devices,
	}, nil
}

// GetDevice returns a specific device by ID
func (c *Controller) GetDevice(ctx context.Context, body *deviceregistrydef.GetDeviceRequest) (*deviceregistrydef.GetDeviceResponse, error) {
	device, err := c.DeviceRepository.Find(body.DeviceId)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find device %q", body.DeviceId)
	}
	if device == nil {
		return nil, oops.NotFound("device %q not found", body.DeviceId)
	}

	// Decorate device with room
	room, err := c.RoomRepository.Find(device.RoomId)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find room %q", device.RoomId)
	}
	if room == nil {
		return nil, oops.NotFound("room %q not found", device.RoomId)
	}
	device.Room = room

	return &deviceregistrydef.GetDeviceResponse{
		DeviceHeader: device,
	}, nil
}
