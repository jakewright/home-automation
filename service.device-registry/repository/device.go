package repository

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	"github.com/jinzhu/copier"

	proto "github.com/jakewright/home-automation/service.device-registry/proto"
)

// DeviceRepository provides access to the underlying storage layer
type DeviceRepository struct {
	// ConfigFilename is the path to the device config file
	ConfigFilename string

	// ReloadInterval is the amount of time to wait before reading from disk again
	ReloadInterval time.Duration

	devices  []*proto.DeviceHeader
	reloaded time.Time
	lock     sync.RWMutex
}

// FindAll returns all devices
func (r *DeviceRepository) FindAll() ([]*proto.DeviceHeader, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var devices []*proto.DeviceHeader
	for _, device := range r.devices {
		out := &proto.DeviceHeader{}
		if err := copier.Copy(&out, device); err != nil {
			return nil, err
		}
		devices = append(devices, out)
	}

	return devices, nil
}

// Find returns a device by ID
func (r *DeviceRepository) Find(id string) (*proto.DeviceHeader, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, device := range r.devices {
		if device.ID == id {
			out := &proto.DeviceHeader{}
			if err := copier.Copy(out, device); err != nil {
				return nil, err
			}

			return out, nil
		}
	}

	return nil, nil
}

// FindByController returns all devices with the given controller name
func (r *DeviceRepository) FindByController(controllerName string) ([]*proto.DeviceHeader, error) {
	// Skip if we've recently reloaded
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var devices []*proto.DeviceHeader
	for _, device := range r.devices {
		if device.ControllerName == controllerName {
			out := &proto.DeviceHeader{}
			if err := copier.Copy(out, device); err != nil {
				return nil, err
			}
			devices = append(devices, out)
		}
	}

	return devices, nil
}

// FindByRoom returns all devices for the given room
func (r *DeviceRepository) FindByRoom(roomID string) ([]*proto.DeviceHeader, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var devices []*proto.DeviceHeader
	for _, device := range r.devices {
		if device.RoomID == roomID {
			out := &proto.DeviceHeader{}
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
		Devices []*proto.DeviceHeader `json:"devices"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.devices = cfg.Devices

	r.reloaded = time.Now()
	return nil
}
