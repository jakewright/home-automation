package repository

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/jinzhu/copier"

	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

// DeviceRepository provides access to the underlying storage layer
type DeviceRepository struct {
	// configFilename is the path to the device config file
	configFilename string

	devices []*deviceregistrydef.DeviceHeader
	lock    sync.RWMutex
}

// NewDeviceRepository returns a new device repository populated with devices
// defined in the JSON file at the given config filename path.
func NewDeviceRepository(configFilename string) (*DeviceRepository, error) {
	r := &DeviceRepository{
		configFilename: configFilename,
	}

	if err := r.reload(); err != nil {
		return nil, err
	}

	return r, nil
}

// FindAll returns all devices
func (r *DeviceRepository) FindAll() ([]*deviceregistrydef.DeviceHeader, error) {

	r.lock.RLock()
	defer r.lock.RUnlock()

	devices := make([]*deviceregistrydef.DeviceHeader, len(r.devices))
	for i, device := range r.devices {
		out := &deviceregistrydef.DeviceHeader{}
		if err := copier.Copy(&out, device); err != nil {
			return nil, err
		}
		devices[i] = out
	}

	return devices, nil
}

// Find returns a device by ID
func (r *DeviceRepository) Find(id string) (*deviceregistrydef.DeviceHeader, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, device := range r.devices {
		if device.Id == id {
			out := &deviceregistrydef.DeviceHeader{}
			if err := copier.Copy(out, device); err != nil {
				return nil, err
			}

			return out, nil
		}
	}

	return nil, nil
}

// FindByController returns all devices with the given controller name
func (r *DeviceRepository) FindByController(controllerName string) ([]*deviceregistrydef.DeviceHeader, error) {
	// Skip if we've recently reloaded
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	// Use an empty slice declaration because if there are
	// no devices we want to make sure an empty list is
	// returned in JSON and not null.
	devices := []*deviceregistrydef.DeviceHeader{}
	for _, device := range r.devices {
		if device.ControllerName == controllerName {
			out := &deviceregistrydef.DeviceHeader{}
			if err := copier.Copy(out, device); err != nil {
				return nil, err
			}
			devices = append(devices, out)
		}
	}

	return devices, nil
}

// FindByRoom returns all devices for the given room
func (r *DeviceRepository) FindByRoom(roomID string) ([]*deviceregistrydef.DeviceHeader, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var devices []*deviceregistrydef.DeviceHeader
	for _, device := range r.devices {
		if device.RoomId == roomID {
			out := &deviceregistrydef.DeviceHeader{}
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
	data, err := ioutil.ReadFile(r.configFilename)
	if err != nil {
		return err
	}

	var cfg struct {
		Devices []*deviceregistrydef.DeviceHeader `json:"devices"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.devices = cfg.Devices
	return nil
}
