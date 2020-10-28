package routes

import "github.com/jakewright/home-automation/services/device-registry/repository"

// Controller handles requests
type Controller struct {
	DeviceRepository *repository.DeviceRepository
	RoomRepository   *repository.RoomRepository
}
