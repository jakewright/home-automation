package main

import (
	"time"

	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.config/controller"
	"github.com/jakewright/home-automation/service.config/domain"
	"github.com/jakewright/home-automation/service.config/service"
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

	r := router.New()
	r.Get("/read/{serviceName}", controller.ReadConfig)
	r.Patch("/reload", controller.ReloadConfig)

	bootstrap.Run(r)
}
