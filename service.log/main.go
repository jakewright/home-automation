package main

import (
	"home-automation/libraries/go/bootstrap"
	"home-automation/libraries/go/router"
	"home-automation/libraries/go/slog"
	"time"

	"home-automation/service.log/dao"

	"home-automation/service.log/routes"
)

func main() {
	if err := bootstrap.Init("service.log"); err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	logRepository := dao.NewLogRepository("/var/log/messages")
	if err := logRepository.Watch(); err != nil {
		slog.Panic("Failed to start watching log file: %v", err)
	}

	controller := routes.Controller{
		Repository: logRepository,
	}

	router.Get("/read", controller.HandleReadLogs)
	router.ListenAndServe()
}
