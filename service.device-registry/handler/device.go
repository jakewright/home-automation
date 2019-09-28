package handler

import (
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/errors"

	"github.com/jakewright/home-automation/service.device-registry/domain"

	"github.com/jakewright/home-automation/libraries/go/request"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.device-registry/repository"
)

// DeviceHandler has functions to handle device-related requests
type DeviceHandler struct {
	DeviceRepository *repository.DeviceRepository
	RoomRepository   *repository.RoomRepository
}

type listRequest struct {
	ControllerName string `json:"controller_name"`
}

type getDeviceRequest struct {
	DeviceID string `json:"device_id"`
}

// HandleListDevices lists all devices known by the registry. Results can be filtered by controller name.
func (h *DeviceHandler) HandleListDevices(w http.ResponseWriter, r *http.Request) {
	body := listRequest{}
	if err := request.Decode(r, &body); err != nil {
		slog.Error("Failed to decode body: %v", err)
		response.WriteJSON(w, err)
		return
	}

	var devices []*domain.Device
	var err error
	if body.ControllerName != "" {
		devices, err = h.DeviceRepository.FindByController(body.ControllerName)
	} else {
		devices, err = h.DeviceRepository.FindAll()
	}
	if err != nil {
		slog.Error("Failed to read devices: %v", err)
		response.WriteJSON(w, err)
		return
	}

	// Decorate the devices with their rooms
	for _, device := range devices {
		room, err := h.RoomRepository.Find(device.RoomID)
		if err != nil {
			slog.Error("Failed to read rooms: %v", err)
			response.WriteJSON(w, err)
			return
		}
		if room == nil {
			slog.Error("Failed to find room %q", device.RoomID)
			response.WriteJSON(w, errors.InternalService("Failed to find room %q", device.RoomID))
			return
		}
		device.Room = room
	}

	response.WriteJSON(w, devices)
}

// HandleGetDevice returns a specific device by ID
func (h *DeviceHandler) HandleGetDevice(w http.ResponseWriter, r *http.Request) {
	body := getDeviceRequest{}
	if err := request.Decode(r, &body); err != nil {
		slog.Error("Failed to decode body: %v", err)
		response.WriteJSON(w, err)
		return
	}

	device, err := h.DeviceRepository.Find(body.DeviceID)
	if err != nil {
		slog.Error("Failed to find device '%s': %v", body.DeviceID, err)
		response.WriteJSON(w, err)
		return
	}

	if device == nil {
		err := errors.NotFound("Device with ID '%s' not found", body.DeviceID)
		response.WriteJSON(w, err)
		return
	}

	// Decorate device with room
	room, err := h.RoomRepository.Find(device.RoomID)
	if err != nil {
		slog.Error("Failed to read rooms: %v", err)
		response.WriteJSON(w, err)
		return
	}

	if room == nil {
		slog.Error("Failed to find room %q", device.RoomID)
		response.WriteJSON(w, errors.InternalService("Failed to find room %q", device.RoomID))
		return
	}

	device.Room = room

	response.WriteJSON(w, device)
}
