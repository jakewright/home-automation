package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/services/dummy/routes"
)

//go:generate jrpc dummy.def

func main() {
	conf := struct{}{}

	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "dummy",
		Config:      &conf,
	})

	routes.Register(svc, &routes.Controller{})

	svc.Run()
}
