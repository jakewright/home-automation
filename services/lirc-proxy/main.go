package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/services/lirc-proxy/routes"
)

//go:generate jrpc lirc_proxy.def

func main() {
	conf := struct{}{}

	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "lirc-proxy",
		Config:      &conf,
	})

	routes.Register(svc, &routes.Controller{})

	svc.Run()
}
