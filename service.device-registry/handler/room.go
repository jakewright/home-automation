package handler

import (
	"github.com/jakewright/home-automation/libraries/go/errors"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"

	"github.com/jakewright/home-automation/service.device-registry/repository"
)

// RoomHandler has functions that deal with room-related requests
type RoomHandler struct {
	DeviceRepository *repository.DeviceRepository
	RoomRepository   *repository.RoomRepository
}

// HandleListRooms returns all rooms known by the registry
func (h *RoomHandler) HandleListRooms(_ *deviceregistrydef.ListRoomsRequest) (*deviceregistrydef.ListRoomsResponse, error) {
	rooms, err := h.RoomRepository.FindAll()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to find rooms")
	}

	// Decorate the rooms with their devices
	for _, room := range rooms {
		devices, err := h.DeviceRepository.FindByRoom(room.Id)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to find devices for message %q", room.Id)
		}
		room.Devices = devices
	}

	return &deviceregistrydef.ListRoomsResponse{
		Rooms: rooms,
	}, nil
}

// HandleGetRoom returns a specific room by ID, including its devices.
func (h *RoomHandler) HandleGetRoom(req *deviceregistrydef.GetRoomRequest) (*deviceregistrydef.GetRoomResponse, error) {
	room, err := h.RoomRepository.Find(req.RoomId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to find room %q", req.RoomId)
	} else if room == nil {
		return nil, errors.NotFound("room %q not found", req.RoomId)
	}

	// Decorate the room with its devices
	devices, err := h.DeviceRepository.FindByRoom(req.RoomId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to find devices for room %q", req.RoomId)
	}
	room.Devices = devices

	return &deviceregistrydef.GetRoomResponse{
		Room: room,
	}, nil
}
