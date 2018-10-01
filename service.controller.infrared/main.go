package main

import (
	"home-automation/libraries/go/bootstrap"
	"log"
)

func main() {
	_, err := bootstrap.NewService("service.controller.infrared")
	if err != nil {
		log.Fatalf("Failed to initialise service: %v", err)
	}
}
