package handler

import "github.com/jakewright/home-automation/service.device-registry/repository"

// Controller handles requests
type Controller struct {
	DeviceRepository *repository.DeviceRepository
	RoomRepository   *repository.RoomRepository
}
