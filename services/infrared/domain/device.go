package domain

import deviceregistrydef "github.com/jakewright/home-automation/services/device-registry/def"

type Device interface {
	ID() string
	Copy() Device
}

func NewDeviceFromDeviceHeader(header *deviceregistrydef.DeviceHeader) (Device, error) {
	// todo
}
