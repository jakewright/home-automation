package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/slog"
	userproto "github.com/jakewright/home-automation/service.user/proto"
)

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.user",
		Database:    true,
	})
	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}

	r := userproto.NewRouter()

	svc.Run(r)
}
