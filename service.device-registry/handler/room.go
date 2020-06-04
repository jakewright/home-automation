package handler

import (
	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

// ListRooms returns all rooms known by the registry
func (c *Controller) ListRooms(r *request, body *deviceregistrydef.ListRoomsRequest) (*deviceregistrydef.ListRoomsResponse, error) {
	rooms, err := c.RoomRepository.FindAll()
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find rooms")
	}

	// Decorate the rooms with their devices
	for _, room := range rooms {
		devices, err := c.DeviceRepository.FindByRoom(room.Id)
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
func (c *Controller) GetRoom(r *request, body *deviceregistrydef.GetRoomRequest) (*deviceregistrydef.GetRoomResponse, error) {
	room, err := c.RoomRepository.Find(body.RoomId)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find room %q", body.RoomId)
	} else if room == nil {
		return nil, oops.NotFound("room %q not found", body.RoomId)
	}

	// Decorate the room with its devices
	devices, err := c.DeviceRepository.FindByRoom(body.RoomId)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to find devices for room %q", body.RoomId)
	}
	room.Devices = devices

	return &deviceregistrydef.GetRoomResponse{
		Room: room,
	}, nil
}
