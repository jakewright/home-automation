package main

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/libraries/go/taxi"
	deviceregistrydef "github.com/jakewright/home-automation/services/device-registry/def"
	"github.com/jakewright/home-automation/services/dmx/dmx"
	"github.com/jakewright/home-automation/services/dmx/domain"
	"github.com/jakewright/home-automation/services/dmx/repository"
	"github.com/jakewright/home-automation/services/dmx/routes"
)

//go:generate jrpc dmx.def

const serviceName = "dmx"

type universeConfig struct {
	UniverseNumber domain.UniverseNumber `envconfig:"UNIVERSE_NUMBER"`
	OLAHost        string                `envconfig:"OLA_HOST"`
	OLAPort        int                   `envconfig:"OLA_PORT"`
}

type config struct {
	Universes []universeConfig `envconfig:"UNIVERSES"`
}

func main() {
	conf := &config{}

	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: serviceName,
		Config:      conf,
	})

	if err := run(svc, conf); err != nil {
		slog.Panicf("Failed to run service: %v", err)
	}
}

func run(svc *bootstrap.Service, conf *config) error {
	client := dmx.NewClient()

	for _, uc := range conf.Universes {
		getSetter, err := dmx.NewOLAClient(uc.OLAHost, uc.OLAPort, uc.UniverseNumber)
		if err != nil {
			return oops.WithMessage(err, "failed to create OLA client")
		}
		client.AddGetSetter(uc.UniverseNumber, getSetter)
	}

	dispatcher := taxi.NewClient()
	repo, err := repository.Init(
		context.Background(),
		serviceName,
		deviceregistrydef.NewClient(dispatcher),
	)
	if err != nil {
		return err
	}

	routes.Register(svc, &routes.Controller{
		Repository: repo,
		Client:     client,
		Publisher:  svc.FirehosePublisher(),
	})

	svc.Run()
	return nil
}
