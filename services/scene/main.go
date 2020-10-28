package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/services/scene/consumer"
	"github.com/jakewright/home-automation/services/scene/routes"
)

//go:generate jrpc scene.def

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.scene",
	})

	firehose.Subscribe(consumer.HandleSetSceneEvent)

	routes.Register(svc, &routes.Controller{
		Database: svc.Database(),
	})

	svc.Run()
}
