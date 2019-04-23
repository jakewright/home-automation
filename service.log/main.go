package main

import (
	"home-automation/libraries/go/bootstrap"
	"home-automation/libraries/go/slog"

	"home-automation/service.log/dao"
	"home-automation/service.log/routes"
)

func main() {
	if err := bootstrap.Init("service.log"); err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	logRepository := dao.NewLogRepository("/var/log/messages")
	go logRepository.Watch()

	(&routes.Controller{
		Repository: logRepository,
	}).RegisterRoutes()

	bootstrap.Run()
}
