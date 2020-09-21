package repository

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/jinzhu/copier"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
)

// DeviceRepository provides access to the underlying storage layer
type DeviceRepository struct {
	// configFilename is the path to the device config file
	configFilename string

	deviceHeaders []*devicedef.Header
	lock          sync.RWMutex
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
func (r *DeviceRepository) FindAll() ([]*devicedef.Header, error) {

	r.lock.RLock()
	defer r.lock.RUnlock()

	devices := make([]*devicedef.Header, len(r.deviceHeaders))
	for i, device := range r.deviceHeaders {
		// todo: since the handler doesn't modify the device
		// anymore (it used to decorate with rooms), there is
		// probably no need to copy the struct.
		out := &devicedef.Header{}
		if err := copier.Copy(&out, device); err != nil {
			return nil, err
		}
		devices[i] = out
	}

	return devices, nil
}

// Find returns a device by ID
func (r *DeviceRepository) Find(id string) (*devicedef.Header, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, device := range r.deviceHeaders {
		if device.GetId() == id {
			out := &devicedef.Header{}
			if err := copier.Copy(out, device); err != nil {
				return nil, err
			}

			return out, nil
		}
	}

	return nil, nil
}

// FindByController returns all devices with the given controller name
func (r *DeviceRepository) FindByController(controllerName string) ([]*devicedef.Header, error) {
	// Skip if we've recently reloaded
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	// Use an empty slice declaration because if there are
	// no devices we want to make sure an empty list is
	// returned in JSON and not null.
	devices := []*devicedef.Header{}
	for _, device := range r.deviceHeaders {
		if device.GetControllerName() == controllerName {
			out := &devicedef.Header{}
			if err := copier.Copy(out, device); err != nil {
				return nil, err
			}
			devices = append(devices, out)
		}
	}

	return devices, nil
}

// FindByRoom returns all devices for the given room
func (r *DeviceRepository) FindByRoom(roomID string) ([]*devicedef.Header, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var devices []*devicedef.Header
	for _, device := range r.deviceHeaders {
		if deviceRoomID, set := device.GetRoomId(); set && deviceRoomID == roomID {
			out := &devicedef.Header{}
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
		Devices []*devicedef.Header `json:"devices"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.deviceHeaders = cfg.Devices
	return nil
}
