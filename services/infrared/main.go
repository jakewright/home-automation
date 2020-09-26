package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/services/infrared/handler"
)

//go:generate jrpc infrared.def

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.infrared",
		Firehose:    true,
	})

	r := handler.NewRouter(svc)

	svc.Run(r)
}
