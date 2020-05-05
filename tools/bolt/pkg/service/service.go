package service

import (
	"fmt"

	"github.com/jakewright/home-automation/libraries/go/util"
	"github.com/jakewright/home-automation/tools/bolt/pkg/compose"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/golang"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

// ListAll returns a list of all home automation services
func ListAll() ([]string, error) {
	var services []string

	systems := getSystems()
	for _, s := range systems {
		ss, err := s.ListAll()
		if err != nil {
			return nil, err
		}

		services = append(services, ss...)
	}

	return util.UniqueStr(services), nil
}

// Run runs the services
func Run(args []string) error {
	names := getServices(args)

	for _, serviceName := range names {
		s, err := getSystem(serviceName)
		if err != nil {
			return fmt.Errorf("failed to get system for %s: %w", serviceName, err)
		}

		if needsBuilding, err := s.NeedsBuilding(serviceName); err != nil {
			return fmt.Errorf("failed to check if needs to build: %w", err)
		} else if needsBuilding {
			output.InfoLn("Building %s...", serviceName)
			if err := s.Build(serviceName); err != nil {
				return fmt.Errorf("failed to build: %w", err)
			}
		}

		op := output.Info("Starting %s", serviceName)
		if err := s.Run(serviceName); err != nil {
			op.Failed()
			return err
		}
		op.Complete()
	}

	return nil
}

// IsRunning returns whether the service is currently running
func IsRunning(name string) (bool, error) {
	s, err := getSystem(name)
	if err != nil {
		return false, err
	}

	return s.IsRunning(name)
}

// Build builds the services
func Build(args []string) error {
	names := getServices(args)

	for _, serviceName := range names {
		s, err := getSystem(serviceName)
		if err != nil {
			return err
		}

		return s.Build(serviceName)
	}

	return nil
}

// Stop stops the services
func Stop(args []string) error {
	names := getServices(args)

	for _, serviceName := range names {
		op := output.Info("Stopping %s", serviceName)
		s, err := getSystem(serviceName)
		if err != nil {
			op.Failed()
			return err
		}

		if err := s.Stop(serviceName); err != nil {
			op.Failed()
			return err
		}
		op.Complete()
	}

	return nil
}

// StopAll stops all services
func StopAll() error {
	output.InfoLn("Stopping all services...")

	for _, s := range getSystems() {
		if err := s.StopAll(); err != nil {
			return err
		}
	}

	return nil
}

// Restart restarts the services
func Restart(args []string) error {
	names := getServices(args)

	for _, serviceName := range names {
		s, err := getSystem(serviceName)
		if err != nil {
			return err
		}

		if running, err := s.IsRunning(serviceName); err != nil {
			return err
		} else if !running {
			output.InfoLn("Cannot restart %s: service not running.", serviceName)
			continue
		}

		op := output.Info("Restarting %s", serviceName)

		if err := s.Stop(serviceName); err != nil {
			op.Failed()
			return err
		}

		if err := s.Run(serviceName); err != nil {
			op.Failed()
			return err
		}
		op.Complete()
	}

	return nil
}

// Exec executes a command inside the service's container
func Exec(serviceName, stdin string, cmd string, args ...string) error {
	s, err := getSystem(serviceName)
	if err != nil {
		return err
	}

	return s.Exec(serviceName, stdin, cmd, args...)
}

type system interface {
	Is(string) (bool, error)
	ListAll() ([]string, error)
	NeedsBuilding(string) (bool, error)
	Run(string) error
	IsRunning(string) (bool, error)
	Build(string) error
	Stop(string) error
	StopAll() error
	Exec(serviceName, stdin string, cmd string, args ...string) error
}

func getSystem(name string) (system, error) {
	for _, s := range getSystems() {
		if is, err := s.Is(name); err != nil {
			return nil, err
		} else if is {
			return s, nil
		}
	}

	return nil, fmt.Errorf("unknown service %q", name)
}

func getSystems() []system {
	return []system{
		// The order is important: compose must be first
		// because golang will "steal" the services.
		&compose.System{},
		&golang.System{},
	}
}

// getServices turns a list of arguments into a
// set of services by expanding the groups.
func getServices(args []string) []string {
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
