package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jakewright/muxinator"
	"home-automation/service.config/controller"
	"home-automation/service.config/domain"
	"home-automation/service.config/service"
)

func main() {
	config := domain.Config{}

	configService := service.ConfigService{
		Location: "/data/config.yaml",
		Config:   &config,
	}

	c := controller.Controller{
		Config:        &config,
		ConfigService: &configService,
	}

	_, err := configService.Reload()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	selfConfig, err := config.Get("service.config")
	if err != nil {
		log.Fatalf("Error reading own config: %v", err)
	}

	if polling, ok := selfConfig["polling"].(map[string]interface{}); ok {
		if enabled, ok := polling["enabled"]; ok && enabled.(bool) {
			if interval, ok := polling["interval"].(int); ok {
				log.Printf("Polling for config changes every %d milliseconds", interval)
				go configService.Watch(time.Millisecond * time.Duration(interval))
			}
		}
	}

	router := muxinator.NewRouter()
	router.Get("/read/{serviceName}", c.ReadConfig)
	router.Patch("/reload", c.ReloadConfig)

	port, ok := selfConfig["port"].(int)
	if !ok {
		port = 80
	}

	log.Printf("Listening on port %d\n", port)
	log.Fatal(router.ListenAndServe(fmt.Sprintf(":%d", port)))
}
