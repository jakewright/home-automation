package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.scene/consumer"
	"github.com/jakewright/home-automation/service.scene/handler"
)

//go:generate jrpc scene.def

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.scene",
		Firehose:    true,
		Database:    true,
	})
	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}

	firehose.Subscribe(consumer.HandleSetSceneEvent)

	r := handler.NewRouter(&handler.Handler{})

	svc.Run(r)
}
