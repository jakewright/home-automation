package service

import (
	"fmt"
	"home-automation/service.config/domain"
	"io/ioutil"
	"log"
	"time"
)

// ConfigService handles loading of the config
type ConfigService struct {
	// Config is a pointer to the config object in use by the service
	Config *domain.Config

	// Location is the path to the YAML file holding the config
	Location string
}

// Watch reads the config file every d duration and applies changes
func (s *ConfigService) Watch(d time.Duration) {
	for {
		_, err := s.Reload()
		if err != nil {
			log.Print(err)
		}

		time.Sleep(d)
	}
}

// Reload reads the config and applies changes
func (s *ConfigService) Reload() (bool, error) {
	data, err := ioutil.ReadFile(s.Location)
	if err != nil {
		return false, err
	}

	reloaded, err := s.Config.SetFromBytes(data)
	if err != nil {
		return false, fmt.Errorf("failed to set config: %v", err)
	}

	if reloaded {
		log.Print("Config reloaded")
	}

	return reloaded, nil
}
