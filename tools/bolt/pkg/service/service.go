package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jakewright/home-automation/tools/bolt/pkg/compose"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/libraries/env"
)

// Expand turns a list of arguments into a
// set of services by expanding the groups.
func Expand(args []string) []string {
	var services []string
	for _, s := range args {
		services = append(services, expandService(s)...)
	}
	return services
}

// expandService returns the set of services
// if s is a group name otherwise s
func expandService(s string) []string {
	for groupName, services := range config.Get().Groups {
		if s == groupName {
			return services
		}
	}

	return []string{s}
}

// Run runs a set of services using the given Compose
func Run(c *compose.Compose, services []string) error {
	for _, s := range services {
		environment, err := getEnv(s)
		if err != nil {
			return fmt.Errorf("failed to get service env: %w", err)
		}

		if err := c.Run(s, environment.AsSh()); err != nil {
			return fmt.Errorf("failed to run service: %w", err)
		}
	}

	return nil
}

func getEnv(service string) (env.Environment, error) {
	var files []string

	configDir := "./private/config/dev/"
	common := filepath.Join(configDir, "common.env")
	serviceSpecific := filepath.Join(configDir, service+".env")

	if exists, err := fileExists(common); err != nil {
		return nil, err
	} else if exists {
		files = append(files, common)
	}

	if exists, err := fileExists(serviceSpecific); err != nil {
		return nil, err
	} else if exists {
		files = append(files, serviceSpecific)
	}

	return env.Parse(files...)
}

func fileExists(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to read %s: %w", name, err)
	}
}
