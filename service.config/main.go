package main

import (
	"home-automation/libraries/go/bootstrap"
	"home-automation/libraries/go/config"
	"home-automation/libraries/go/router"
	"home-automation/libraries/go/slog"
	"time"

	"home-automation/service.config/controller"
	"home-automation/service.config/domain"
	"home-automation/service.config/service"
)

func main() {
	c := domain.Config{}

	configService := service.ConfigService{
		Location: "/data/config.yaml",
		Config:   &c,
	}

	controller := controller.Controller{
		Config:        &c,
		ConfigService: &configService,
	}

	_, err := configService.Reload()
	if err != nil {
		slog.Panic("Failed to load config: %v", err)
	}

	selfConfig, err := c.Get("service.config")
	if err != nil {
		slog.Panic("Error reading own config: %v", err)
	}

	config.DefaultProvider = config.New(selfConfig)

	if config.Get("polling.enabled").Bool(false) {
		interval := config.Get("polling.interval").Int(30000)
		slog.Info("Polling for config changes every %d milliseconds", interval)
		go configService.Watch(time.Millisecond * time.Duration(interval))
	}

	router.Get("/read/{serviceName}", controller.ReadConfig)
	router.Patch("/reload", controller.ReloadConfig)

	bootstrap.Run()
}
