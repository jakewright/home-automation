package service

import (
	"fmt"
	"strings"

	"github.com/jakewright/home-automation/tools/bolt/pkg/compose"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
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

	// To allow the user to specify the directory of the service instead of
	// just the service name, strip the potential prefixes off the string.
	s = strings.TrimPrefix(s, "service/")
	s = strings.TrimPrefix(s, "./service/")

	return []string{s}
}

// Run runs a set of services using the given Compose
func Run(c *compose.Compose, services []string) error {
	for _, s := range services {
		if err := c.Run(s); err != nil {
			return fmt.Errorf("failed to run service: %w", err)
		}
	}

	return nil
}
