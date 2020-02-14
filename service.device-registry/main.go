package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.device-registry/handler"
	"github.com/jakewright/home-automation/service.device-registry/repository"
)

//go:generate jrpc deviceregistry.def

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.device-registry",
	})

	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}

	configFilename := config.Get("configFilename").String()
	reloadInterval := config.Get("reloadInterval").Duration()

	if configFilename == "" {
		slog.Panicf("configFilename is empty")
	}

	if reloadInterval == 0 {
		slog.Panicf("reloadInterval is empty")
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

	svc.Run(r)
}
