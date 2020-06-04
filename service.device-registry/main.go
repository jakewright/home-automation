package main

import (
	"time"

	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.device-registry/handler"
	"github.com/jakewright/home-automation/service.device-registry/repository"
)

//go:generate jrpc deviceregistry.def

func main() {
	conf := struct {
		ConfigFilename string
		ReloadInterval time.Duration
	}{}

	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.device-registry",
		Config:      &conf,
	})

	if conf.ConfigFilename == "" {
		slog.Panicf("configFilename is empty")
	}

	if conf.ReloadInterval == 0 {
		slog.Panicf("reloadInterval is empty")
	}

	dr := &repository.DeviceRepository{
		ConfigFilename: conf.ConfigFilename,
		ReloadInterval: conf.ReloadInterval,
	}
	rr := &repository.RoomRepository{
		ConfigFilename: conf.ConfigFilename,
		ReloadInterval: conf.ReloadInterval,
	}

	r := handler.NewRouter(&handler.Controller{
		DeviceRepository: dr,
		RoomRepository:   rr,
	})

	svc.Run(r)
}
