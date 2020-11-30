package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/services/infrared/routes"
)

//go:generate jrpc infrared.def

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.infrared",
		Firehose:    true,
	})

	r := routes.Register(svc, &routes.Controller{
		Repository: nil,
		IR:         nil,
	})

	svc.Run(r)
}
