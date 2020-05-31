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

	h := &handler.Handler{
		TemplateDirectory: conf.TemplateDirectory,
		LogRepository:     logRepository,
		Watcher:           watcher,
	}

	r := router.New()
	r.Get("/", h.HandleRead)
	r.Get("/ws", h.HandleWebSocket)
	r.Post("/write", h.HandleWrite)

	svc.Run(r, watcher)
}
