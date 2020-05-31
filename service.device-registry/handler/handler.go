package handler

import "github.com/jakewright/home-automation/service.device-registry/repository"

// Handler handles requests
type Handler struct {
	DeviceRepository *repository.DeviceRepository
	RoomRepository   *repository.RoomRepository
}
