package repository

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	"github.com/jakewright/home-automation/service.device-registry/domain"
)

type RoomRepository struct {
	// ConfigFilename is the path to the room config file
	ConfigFilename string

	// ReloadInterval is the amount of time to wait before reading from disk again
	ReloadInterval time.Duration

	rooms    map[string]*domain.Room
	reloaded time.Time
	lock     sync.RWMutex
}

func (r *RoomRepository) FindAll() ([]*domain.Room, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var rooms []*domain.Room
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepository) Find(id string) (*domain.Room, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	return r.rooms[id], nil
}

func (r *RoomRepository) reload() error {
	// Skip if we've recently reloaded
	if r.reloaded.Add(r.ReloadInterval).After(time.Now()) {
		return nil
	}

	data, err := ioutil.ReadFile(r.ConfigFilename)
	if err != nil {
		return err
	}

	var cfg struct {
		Rooms []*domain.Room `json:"rooms"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if r.rooms == nil {
		r.rooms = map[string]*domain.Room{}
	}

	for _, room := range cfg.Rooms {
		r.rooms[room.ID] = room
	}

	r.reloaded = time.Now()
	return nil
}
