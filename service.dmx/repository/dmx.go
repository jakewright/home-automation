package repository

import (
	"fmt"

	"github.com/jakewright/home-automation/libraries/go/rpc"
)

type DMXRepository struct {
	ServiceName string
}

func (r *DMXRepository) FetchDevices() {
	url := fmt.Sprintf("service.device-registry/devices?controller_name=%s", r.ServiceName)
	if _, err := rpc.Get(url); err != nil {

	}
}
