package main

import "github.com/jakewright/home-automation/libraries/go/bootstrap"

const serviceName = "fake-dmx"

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: serviceName,
	})

	svc.Run()
}
