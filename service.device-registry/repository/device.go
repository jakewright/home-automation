package repository

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	"github.com/jinzhu/copier"

	"github.com/jakewright/home-automation/service.device-registry/domain"
)

type DeviceRepository struct {
	// ConfigFilename is the path to the device config file
	ConfigFilename string

	// ReloadInterval is the amount of time to wait before reading from disk again
	ReloadInterval time.Duration

	devices  map[string]*domain.Device
	reloaded time.Time
	lock     sync.RWMutex
}

func (r *DeviceRepository) FindAll() ([]*domain.Device, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var devices []*domain.Device
	for _, device := range r.devices {
		out := &domain.Device{}
		if err := copier.Copy(&out, device); err != nil {
			return nil, err
		}
		devices = append(devices, out)
	}

	return devices, nil
}

func (r *DeviceRepository) Find(deviceID string) (*domain.Device, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	out := &domain.Device{}
	if err := copier.Copy(out, r.devices[deviceID]); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *DeviceRepository) FindByController(controllerName string) ([]*domain.Device, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var devices []*domain.Device
	for _, device := range r.devices {
		if device.ControllerName == controllerName {
			out := &domain.Device{}
			if err := copier.Copy(out, device); err != nil {
				return nil, err
			}
			devices = append(devices, out)
		}
	}

	return devices, nil
}

func (r *DeviceRepository) FindByRoom(roomID string) ([]*domain.Device, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var devices []*domain.Device
	for _, device := range r.devices {
		if device.RoomID == roomID {
			out := &domain.Device{}
			if err := copier.Copy(out, device); err != nil {
				return nil, err
			}
			devices = append(devices, out)
		}
	}

	return devices, nil
}

// reload reads the config and applies changes
func (r *DeviceRepository) reload() error {
	// Skip if we've recently reloaded
	if r.reloaded.Add(r.ReloadInterval).After(time.Now()) {
		return nil
	}

	data, err := ioutil.ReadFile(r.ConfigFilename)
	if err != nil {
		return err
	}

	var cfg struct {
		Devices []*domain.Device `json:"devices"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if r.devices == nil {
		r.devices = map[string]*domain.Device{}
	}

	for _, device := range cfg.Devices {
		r.devices[device.ID] = device
	}

	r.reloaded = time.Now()
	return nil
}
