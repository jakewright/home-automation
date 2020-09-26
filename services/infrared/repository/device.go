package repository

import (
	"sync"

	"github.com/jakewright/home-automation/services/infrared/domain"
)

// DeviceRepository holds devices
type DeviceRepository struct {
	devices map[string]domain.Device
	mux     *sync.RWMutex
}

// New returns a new DeviceRepository
func New() *DeviceRepository {
	return &DeviceRepository{
		devices: make(map[string]domain.Device),
		mux:     &sync.RWMutex{},
	}
}

// Find returns the device with the given ID
func (r *DeviceRepository) Find(id string) domain.Device {
	r.mux.RLock()
	defer r.mux.RUnlock()

	d, ok := r.devices[id]
	if !ok {
		return nil
	}

	return d.Copy()
}

// AddDevice adds the given device to the
// repository if it does not already exist
func (r *DeviceRepository) AddDevice(d domain.Device) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.devices[d.ID()]; ok {
		return
	}

	r.devices[d.ID()] = d
}

// Save adds the given device to the repository
// replacing any existing device with the same ID
func (r *DeviceRepository) Save(d domain.Device) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.devices[d.ID()] = d
}
