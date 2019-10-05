package repository

import (
	"fmt"
	"sync"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/service.dmx/domain"

	"github.com/jakewright/home-automation/libraries/go/rpc"
)

type DMXRepository struct {
	ServiceName string
	universe    *domain.Universe
	mux         sync.RWMutex
}

func (r *DMXRepository) Find(id string) domain.Fixture {
	r.mux.RLock()
	defer r.mux.RUnlock()

	for _, f := range r.universe.Fixtures {
		if f.ID() == id {
			return f
		}
	}

	return nil
}

func (r *DMXRepository) FetchDevices() error {
	url := fmt.Sprintf("service.device-registry/devices?controller_name=%s", r.ServiceName)
	var rsp []*domain.DeviceHeader
	if _, err := rpc.Get(url, &rsp); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	for _, device := range rsp {
		switch {
		case device.ControllerName != r.ServiceName:
			return errors.InternalService("Device %s is not for this controller", device.ID)
		case device.Type != "dmx":
			return errors.InternalService("Device %s does not have type dmx", device.ID)
		}

		fixture, err := domain.NewFixtureFromDeviceHeader(device)
		if err != nil {
			return errors.InternalService("Failed to create fixture: %v", err)
		}
		r.universe.Fixtures = append(r.universe.Fixtures, fixture)
	}
}
