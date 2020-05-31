package handler

import (
	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

// ListRooms returns all rooms known by the registry
func (h *Handler) ListRooms(r *Request, body *deviceregistrydef.ListRoomsRequest) (*deviceregistrydef.ListRoomsResponse, error) {
	rooms, err := h.RoomRepository.FindAll()
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find rooms")
	}

	// Decorate the rooms with their devices
	for _, room := range rooms {
		devices, err := h.DeviceRepository.FindByRoom(room.Id)
		if err != nil {
			return nil, oops.WithMessage(err, "failed to find devices for message %q", room.Id)
		}
		room.Devices = devices
	}

	return &deviceregistrydef.ListRoomsResponse{
		Rooms: rooms,
	}, nil
}

// GetRoom returns a specific room by ID, including its devices.
func (h *Handler) GetRoom(r *Request, body *deviceregistrydef.GetRoomRequest) (*deviceregistrydef.GetRoomResponse, error) {
	room, err := h.RoomRepository.Find(body.RoomId)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find room %q", body.RoomId)
	} else if room == nil {
		return nil, oops.NotFound("room %q not found", body.RoomId)
	}

	// Decorate the room with its devices
	devices, err := h.DeviceRepository.FindByRoom(body.RoomId)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find devices for room %q", body.RoomId)
	}
	room.Devices = devices

	return &deviceregistrydef.GetRoomResponse{
		Room: room,
	}, nil
}
