package handler

import (
	"github.com/jakewright/home-automation/libraries/go/errors"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"

	"github.com/jakewright/home-automation/service.device-registry/repository"
)

// DeviceHandler has functions to handle device-related requests
type DeviceHandler struct {
	DeviceRepository *repository.DeviceRepository
	RoomRepository   *repository.RoomRepository
}

// HandleListDevices lists all devices known by the registry. Results can be filtered by controller name.
func (h *DeviceHandler) HandleListDevices(req *deviceregistrydef.ListDevicesRequest) (*deviceregistrydef.ListDevicesResponse, error) {
	var devices []*deviceregistrydef.DeviceHeader
	var err error
	if req.ControllerName != "" {
		devices, err = h.DeviceRepository.FindByController(req.ControllerName)
	} else {
		devices, err = h.DeviceRepository.FindAll()
	}
	if err != nil {
		return nil, errors.WithMessage(err, "failed to find devices")
	}

	// Decorate the devices with their rooms
	for _, device := range devices {
		room, err := h.RoomRepository.Find(device.RoomId)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to find room %q", device.RoomId)
		}
		if room == nil {
			return nil, errors.NotFound("room %q not found", device.RoomId)
		}

		device.Room = room
	}

	return &deviceregistrydef.ListDevicesResponse{
		DeviceHeaders: devices,
	}, nil
}

// HandleGetDevice returns a specific device by ID
func (h *DeviceHandler) HandleGetDevice(req *deviceregistrydef.GetDeviceRequest) (*deviceregistrydef.GetDeviceResponse, error) {
	device, err := h.DeviceRepository.Find(req.DeviceId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to find device %q", req.DeviceId)
	}
	if device == nil {
		return nil, errors.NotFound("device %q not found", req.DeviceId)
	}

	// Decorate device with room
	room, err := h.RoomRepository.Find(device.RoomId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to find room %q", device.RoomId)
	}
	if room == nil {
		return nil, errors.NotFound("room %q not found", device.RoomId)
	}
	device.Room = room

	return &deviceregistrydef.GetDeviceResponse{
		DeviceHeader: device,
	}, nil
}
