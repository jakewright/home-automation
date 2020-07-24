package compose

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/util"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/docker"
)

// Compose performs docker-compose tasks
type Compose struct {
	f *composeFile
}

// New returns a new Compose struct
func New() (*Compose, error) {
	f, err := parse(config.Get().DockerComposeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse docker-compose file: %w", err)
	}

	return &Compose{f}, nil
}

// ListAll returns a list of all service names
func (c *Compose) ListAll() ([]string, error) {
	services := make([]string, 0, len(c.f.Services))
	for name := range c.f.Services {
		services = append(services, name)
	}
	return services, nil
}

// Run starts the service, building first if necessary.
func (c *Compose) Run(service string) error {
	args := []string{"up", "-d", "--renew-anon-volumes", "--remove-orphans", service}
	if err := c.cmd(args...).SetPseudoTTY().Run().Err; err != nil {
		return fmt.Errorf("failed to cmd docker-compose up: %w", err)
	}

	return nil
}

// IsRunning returns whether the service is currently running
func (c *Compose) IsRunning(serviceName string) (bool, error) {
	containerID, err := c.getContainerID(serviceName)
	if err != nil {
		return false, err
	} else if containerID == "" {
		return false, nil
	}

	return docker.IsContainerRunning(containerID)
}

// Build builds the service
func (c *Compose) Build(services []string) error {
	args := append([]string{"build", "--pull"}, services...)
	if err := c.cmd(args...).SetPseudoTTY().Run().Err; err != nil {
		return err
	}

	return nil
}

// Stop stops the service
func (c *Compose) Stop(services []string) error {
	args := append([]string{"stop"}, services...)
	if err := c.cmd(args...).SetPseudoTTY().Run().Err; err != nil {
		return fmt.Errorf("failed to cmd docker-compose stop: %w", err)
	}

	return nil

}

// StopAll stops all docker-compose services
func (c *Compose) StopAll() error {
	if err := c.cmd("stop").SetPseudoTTY().Run().Err; err != nil {
		return fmt.Errorf("failed to cmd docker-compose stop: %w", err)
	}

	return nil
}

// Exec executes the command inside the service's container
func (c *Compose) Exec(serviceName, stdin string, cmd string, args ...string) error {
	containerID, err := c.getContainerID(serviceName)
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

// Ports returns the host ports that the service exposes
func (c *Compose) Ports(serviceName string) ([]string, error) {
	raw := c.f.Services[serviceName].Ports
	var ports []string

	for _, port := range raw {
		// This only supports the 3000:3000 syntax
		re := regexp.MustCompile(`^(\d+):(\d+)$`)
		matches := re.FindStringSubmatch(port)
		if len(matches) != 3 {
			return nil, fmt.Errorf("unsupported port format: %s", port)
		}

		ports = append(ports, matches[1])
	}

	return ports, nil
}

// Logs runs docker-compose logs
func (c *Compose) Logs(services []string) error {
	args := []string{"logs", "--follow", "--timestamps", "--tail=30"}
	args = append(args, services...)
	// Ignore the error because it returns a non-zero exit code on Ctrl+c
	c.cmd(args...).SetPseudoTTY().Run()
	return nil
}

// cmd returns a docker-compose command with the
// docker-compose file and project name flags set
func (c *Compose) cmd(args ...string) *exe.Cmd {
	a := []string{"-f", config.Get().DockerComposeFilePath, "-p", config.Get().ProjectName}
	a = append(a, args...)
	return exe.Command("docker-compose", a...)
}

func (c *Compose) getContainerID(serviceName string) (string, error) {
	result := c.cmd("ps", "-q", serviceName).Run()
	if result.Err != nil {
		return "", fmt.Errorf("failed to cmd docker-compose ps: %w", result.Err)
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
