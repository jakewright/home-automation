package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/services/user/handler"
)

//go:generate jrpc user.def

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.user",
		Database:    true,
	})

	r := handler.NewRouter(&handler.Controller{})
	svc.Run(r)
}
