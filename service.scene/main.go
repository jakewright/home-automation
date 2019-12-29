package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName:"service.scene",
		Database: true,
	})
	if err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	r := router.New()

	svc.Run(r)
}
