package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.user/handler"
)

//go:generate jrpc user.def

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.user",
		Database:    true,
	})
	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}

	r := handler.NewRouter(&handler.Handler{})
	svc.Run(r)
}
