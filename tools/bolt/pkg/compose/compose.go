package compose

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/util"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/docker"
)

// System is the docker-compose service management system
type System struct {
	f *composeFile
}

// Is returns whether the given service is defined in the docker-compose.yml composeFile
func (s *System) Is(serviceName string) (bool, error) {
	f, err := s.file()
	if err != nil {
		return false, err
	}

	for composeService := range f.Services {
		if serviceName == composeService {
			return true, nil
		}
	}

	return false, nil
}

// ListAll returns a list of all service names
func (s *System) ListAll() ([]string, error) {
	f, err := s.file()
	if err != nil {
		return nil, err
	}

	services := make([]string, 0, len(f.Services))
	for name := range f.Services {
		services = append(services, name)
	}
	return services, nil
}

// NeedsBuilding returns whether the service needs to be built before it can be run
func (s *System) NeedsBuilding(serviceName string) (bool, error) {
	f, err := s.file()
	if err != nil {
		return false, err
	}

	// If there's an image name defined in the compose
	// file, we can check for its existence.
	if f.Services[serviceName].Image != "" {
		exists, err := docker.ImageExists(f.Services[serviceName].Image)
		if err != nil {
			return false, fmt.Errorf("failed to check if image exists: %w", err)
		}

		return !exists, nil
	}

	// If we don't know the image name, just check whether the container exists.
	// Obviously the image could exist without the container but meh, close enough.
	containerID, err := s.getContainerID(serviceName)
	if err != nil {
		return false, fmt.Errorf("failed to get container ID: %w", err)
	}

	return containerID == "", nil
}

// Run starts the service, building first if necessary.
func (s *System) Run(serviceName string) error {
	args := []string{"up", "-d", "--renew-anon-volumes", "--remove-orphans", serviceName}

	if err := s.dockerCompose(args...).Run().Err; err != nil {
		return fmt.Errorf("failed to run docker-compose up: %w", err)
	}

	return nil
}

// IsRunning returns whether the service is currently running
func (s *System) IsRunning(serviceName string) (bool, error) {
	containerID, err := s.getContainerID(serviceName)
	if err != nil {
		return false, err
	} else if containerID == "" {
		return false, nil
	}

	return docker.IsContainerRunning(containerID)
}

// Build builds the service
func (s *System) Build(serviceName string) error {
	args := []string{"build", "--pull", serviceName}
	if err := s.dockerCompose(args...).SetPseudoTTY().Run().Err; err != nil {
		return err
	}
	return nil
}

// Stop stops the service
func (s *System) Stop(serviceName string) error {
	if err := s.dockerCompose("stop", serviceName).Run().Err; err != nil {
		return fmt.Errorf("failed to run docker-compose stop: %w", err)
	}

	return nil

}

// StopAll stops all docker-compose services
func (s *System) StopAll() error {
	// The output of this doesn't match the rest of the tool but it's
	// too much effort to get the list of running containers and stop
	// each one manually.
	if err := s.dockerCompose("stop").SetPseudoTTY().Run().Err; err != nil {
		return fmt.Errorf("failed to run docker-compose stop: %w", err)
	}

	return nil
}

// Exec executes the command inside the service's container
func (s *System) Exec(serviceName, stdin string, cmd string, args ...string) error {
	containerID, err := s.getContainerID(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get container ID of %s: %w", serviceName, err)
	}

	joinedCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	dockerArgs := []string{"exec", "-i", containerID, "sh", "-c", joinedCmd}

	if err := exe.Command("docker", dockerArgs...).SetInput(stdin).Run().Err; err != nil {
		return fmt.Errorf("failed to docker exec: %w", err)
	}

	return nil
}

type composeFile struct {
	Version  string                     `yaml:"version"`
	Services map[string]*composeService `yaml:"services"`
	Networks map[string]interface{}     `yaml:"networks"`
}

type composeService struct {
	Image string `yaml:"image"`
}

func (s *System) file() (*composeFile, error) {
	if s.f == nil {
		b, err := ioutil.ReadFile(config.Get().DockerComposeFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read docker-compose composeFile: %w", err)
		}

		s.f = &composeFile{}
		if err := yaml.Unmarshal(b, s.f); err != nil {
			return nil, fmt.Errorf("failed to unmarshal docker-compose composeFile: %w", err)
		}

		if len(s.f.Networks) > 0 {
			return nil, fmt.Errorf("custom docker-compose networks are not supported")
		}
	}

	return s.f, nil
}

// dockerCompose returns a docker-compose command with the
// docker-compose file and project name flags set
func (s *System) dockerCompose(args ...string) *exe.Cmd {
	a := []string{"-f", config.Get().DockerComposeFilePath, "-p", config.Get().ProjectName}
	a = append(a, args...)
	return exe.Command("docker-compose", a...)
}

func (s *System) getContainerID(serviceName string) (string, error) {
	result := s.dockerCompose("ps", "-q", serviceName).Run()
	if result.Err != nil {
		return "", fmt.Errorf("failed to run docker-compose ps: %w", result.Err)
	}

	lines := strings.Split(result.Stdout, "\n")
	lines = util.RemoveWhitespaceStrings(lines)

	if len(lines) == 0 {
		return "", nil
	} else if len(lines) > 1 {
		return "", fmt.Errorf("found multiple containers")
	}

	return lines[0], nil
}
