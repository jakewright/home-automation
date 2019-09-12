package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.device-registry/handler"
	"github.com/jakewright/home-automation/service.device-registry/repository"
)

func main() {
	if err := bootstrap.Init("service.device-registry"); err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	configFilename := config.Get("configFilename").String()
	reloadInterval := config.Get("reloadInterval").Duration()

	if configFilename == "" {
		slog.Panic("configFilename is empty")
	}

	if reloadInterval == 0 {
		slog.Panic("reloadInterval is empty")
	}

	dr := &repository.DeviceRepository{
		ConfigFilename: configFilename,
		ReloadInterval: reloadInterval,
	}
	rr := &repository.RoomRepository{
		ConfigFilename: configFilename,
		ReloadInterval: reloadInterval,
	}

	deviceHandler := handler.DeviceHandler{
		DeviceRepository: dr,
		RoomRepository:   rr,
	}
	roomHandler := handler.RoomHandler{
		DeviceRepository: dr,
		RoomRepository:   rr,
	}

	r := router.New()
	r.Get("/devices", deviceHandler.HandleListDevices)
	r.Get("/device/{device_id}", deviceHandler.HandleGetDevice)
	r.Get("/rooms", roomHandler.HandleListRooms)
	r.Get("/room/{room_id}", roomHandler.HandleGetRoom)

	bootstrap.Run(r)
}
