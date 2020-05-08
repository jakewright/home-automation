package service

import "github.com/jakewright/home-automation/tools/bolt/pkg/config"

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
