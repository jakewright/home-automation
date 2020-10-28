package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/services/user/routes"
)

//go:generate jrpc user.def

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.user",
	})

	routes.Register(svc, &routes.Controller{})

	svc.Run()
}
