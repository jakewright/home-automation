package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/services/device-registry/handler"
	"github.com/jakewright/home-automation/services/device-registry/repository"
)

//go:generate jrpc deviceregistry.def

func main() {
	conf := struct {
		ConfigFilename string
	}{}

	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.device-registry",
		Config:      &conf,
	})

	if conf.ConfigFilename == "" {
		slog.Panicf("configFilename is empty")
	}

	dr, err := repository.NewDeviceRepository(conf.ConfigFilename)
	if err != nil {
		slog.Panicf("failed to init device repository: %v", err)
	}

	rr, err := repository.NewRoomRepository(conf.ConfigFilename)
	if err != nil {
		slog.Panicf("failed to init room repository: %v", err)
	}

	r := handler.NewRouter(&handler.Controller{
		DeviceRepository: dr,
		RoomRepository:   rr,
	})

	svc.Run(r)
}
