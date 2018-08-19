package main

import (
	"home-automation/libraries/go/muxinator"
	"home-automation/service.config/controller"
	"home-automation/service.config/domain"
	"home-automation/service.config/service"
	"log"
	"net/http"
	"time"
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

	go configService.Watch(time.Second * 30)

	router := muxinator.NewRouter()
	router.Get("/read/{serviceName}", c.ReadConfig)
	router.Patch("/reload", c.ReloadConfig)

	log.Println("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", router.BuildHandler()))
}
