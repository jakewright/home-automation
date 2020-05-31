package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.log/handler"
	"github.com/jakewright/home-automation/service.log/repository"
	"github.com/jakewright/home-automation/service.log/watch"
)

func main() {
	conf := struct {
		LogDirectory      string
		TemplateDirectory string
	}{}

	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.log",
		Config:      &conf,
	})

	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}

	if conf.LogDirectory == "" {
		slog.Panicf("logDirectory not set in config")
	}

	if conf.TemplateDirectory == "" {
		slog.Panicf("templateDirectory not set in config")
	}

	logRepository := &repository.LogRepository{
		LogDirectory: conf.LogDirectory,
	}

	watcher := &watch.Watcher{
		LogRepository: logRepository,
	}

	readHandler := handler.ReadHandler{
		TemplateDirectory: conf.TemplateDirectory,
		LogRepository:     logRepository,
		Watcher:           watcher,
	}

	r := router.New()
	r.Get("/", readHandler.HandleRead)
	r.Get("/ws", readHandler.HandleWebSocket)
	r.Post("/write", handler.HandleWrite)

	svc.Run(r, watcher)
}
