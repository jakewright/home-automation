package handler

import (
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/errors"

	"github.com/jakewright/home-automation/libraries/go/request"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.device-registry/repository"
)

type RoomHandler struct {
	DeviceRepository *repository.DeviceRepository
	RoomRepository   *repository.RoomRepository
}

type getRoomRequest struct {
	RoomID string `json:"room_id"`
}

func (h *RoomHandler) HandleListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.RoomRepository.FindAll()
	if err != nil {
		slog.Error("Failed to find rooms: %v", err)
		response.WriteJSON(w, err)
		return
	}

	// Decorate the rooms with their devices
	for _, room := range rooms {
		devices, err := h.DeviceRepository.FindByRoom(room.ID)
		if err != nil {
			slog.Error("Failed to find devices for room '%s': %v", room.ID, err)
			response.WriteJSON(w, err)
			return
		}
		room.Devices = devices
	}

	response.WriteJSON(w, rooms)
}

func (h *RoomHandler) HandleGetRoom(w http.ResponseWriter, r *http.Request) {
	body := getRoomRequest{}
	if err := request.Decode(r, &body); err != nil {
		slog.Error("Failed to decode body: %v", err)
		response.WriteJSON(w, err)
		return
	}

	room, err := h.RoomRepository.Find(body.RoomID)
	if err != nil {
		slog.Error("Failed to find room '%s': %v", body.RoomID, err)
		response.WriteJSON(w, err)
		return
	}

	if room == nil {
		err := errors.NotFound("Room with ID '%s' not found", body.RoomID)
		response.WriteJSON(w, err)
		return
	}

	// Decorate the room with its devices
	devices, err := h.DeviceRepository.FindByRoom(body.RoomID)
	if err != nil {
		slog.Error("Failed to find devices for room '%s': %v", body.RoomID, err)
		response.WriteJSON(w, err)
		return
	}

	room.Devices = devices
	response.WriteJSON(w, room)
}
