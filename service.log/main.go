package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.log/handler"
	"github.com/jakewright/home-automation/service.log/repository"
	"github.com/jakewright/home-automation/service.log/watch"
)

func main() {
	if err := bootstrap.Init("service.log"); err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	logDirectory := config.Get("logDirectory").String()
	if logDirectory == "" {
		slog.Panic("Log directory not set in config")
	}

	logRepository := &repository.LogRepository{
		LogDirectory: logDirectory,
	}

	watcher := &watch.Watcher{
		LogRepository: logRepository,
	}

	readHandler := handler.ReadHandler{
		LogRepository: logRepository,
		Watcher:       watcher,
	}

	r := router.New()
	r.Get("/", readHandler.HandleRead, readHandler.DecodeBody)
	r.Get("/ws", readHandler.HandleWebSocket, readHandler.DecodeBody)
	r.Post("/write", handler.HandleWrite)

	bootstrap.Run(r, watcher)
}
