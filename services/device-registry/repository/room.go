package repository

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/jinzhu/copier"

	deviceregistrydef "github.com/jakewright/home-automation/services/device-registry/def"
)

// RoomRepository provides access to the underlying storage layer
type RoomRepository struct {
	// configFilename is the path to the room config file
	configFilename string

	rooms []*deviceregistrydef.Room
	lock  sync.RWMutex
}

// NewRoomRepository returns a new room repository populated with rooms defined
// in the JSON file at the given config filename path.
func NewRoomRepository(configFilename string) (*RoomRepository, error) {
	r := &RoomRepository{
		configFilename: configFilename,
	}

	if err := r.reload(); err != nil {
		return nil, err
	}

	return r, nil
}

// FindAll returns all rooms
func (r *RoomRepository) FindAll() ([]*deviceregistrydef.Room, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var rooms []*deviceregistrydef.Room
	for _, room := range r.rooms {
		out := &deviceregistrydef.Room{}
		if err := copier.Copy(out, room); err != nil {
			return nil, err
		}
		rooms = append(rooms, out)
	}

	return rooms, nil
}

// Find returns a room by ID
func (r *RoomRepository) Find(id string) (*deviceregistrydef.Room, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, room := range r.rooms {
		if room.GetId() == id {
			out := &deviceregistrydef.Room{}
			if err := copier.Copy(out, room); err != nil {
				return nil, err
			}

			return out, nil
		}
	}

	return nil, nil
}

func (r *RoomRepository) reload() error {
	data, err := ioutil.ReadFile(r.configFilename)
	if err != nil {
		return err
	}

	var cfg struct {
		Rooms []*deviceregistrydef.Room `json:"rooms"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.rooms = cfg.Rooms
	return nil
}
