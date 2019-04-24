package main

import (
	"home-automation/libraries/go/bootstrap"
	"home-automation/libraries/go/router"
	"home-automation/libraries/go/slog"

	"home-automation/service.log/handler"
	"home-automation/service.log/repository"
	"home-automation/service.log/watch"
)

func main() {
	if err := bootstrap.Init("service.log"); err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	fileLocation := "/var/log/messages"

	logRepository := &repository.LogRepository{
		Location: fileLocation,
	}

	watcher := &watch.Watcher{
		LogRepository: logRepository,
		Location:      fileLocation,
	}

	h := handler.LogHandler{
		LogRepository: logRepository,
		Watcher:       watcher,
	}

	r := router.New()
	r.Get("/", h.HandleRead, h.DecodeBody)
	r.Get("/ws", h.HandleWebSocket, h.DecodeBody)

	bootstrap.Run(r, watcher)
}
