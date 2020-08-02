package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/services/scene/consumer"
	"github.com/jakewright/home-automation/services/scene/handler"
)

//go:generate jrpc scene.def

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.scene",
		Firehose:    true,
		Database:    true,
	})

	firehose.Subscribe(consumer.HandleSetSceneEvent)

	r := handler.NewRouter(svc, &handler.Controller{})

	svc.Run(r)
}
