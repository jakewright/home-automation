package docker

import (
	"fmt"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
)

const runSchemeLabel = "home_automation.run_scheme"

// ImageForService returns whether an image for the service is already built
func ImageForService(serviceName string) (bool, error) {
	return ImageExists(imageName(serviceName))
}

func containerID(serviceName string) (string, error) {
	result := exe.Command("docker", "ps", "-q", "-f", "name="+containerName(serviceName)).Run()
	if result.Err != nil {
		return "", fmt.Errorf("failed to run docker ps: %w", result.Err)
	}

	return result.Stdout, nil
}

// IsRunning returns whether a service is running
func IsRunning(serviceName string) (bool, error) {
	id, err := containerID(serviceName)
	if err != nil {
		return false, err
	} else if id == "" {
		return false, nil
	}

	return IsContainerRunning(id)
}

// Run runs the service. It assumes the image has already been built.
func Run(serviceName, runScheme string) error {
	// Run in detached mode and remove the container after exit
	args := []string{"run", "-d", "--rm"}

	// Name is used to check whether a container is already running
	args = append(args, "--name", containerName(serviceName))

	// Join the docker-compose network
	args = append(args, "--network", networkName())

	// Make this service discoverable to other containers by its service name.
	// The docs say this should be a list but it only seems to work like this.
	args = append(args, "--network-alias", serviceName)

	// Bind each exposed port to a random port on the hose
	args = append(args, "-P")

	// Set a label so the calling scheme (e.g. golang) can find this service again
	args = append(args, "--label", fmt.Sprintf("%s=%s", runSchemeLabel, runScheme))

	args = append(args, imageName(serviceName))

	// Run without a pty because all it outputs is the container ID
	result := exe.Command("docker", args...).Run()
	if result.Err != nil {
		return fmt.Errorf("failed to run: %w", result.Err)
	}

	return nil
}

// Build builds a docker image for the service using the given docker file
func Build(serviceName, dockerfileFilename string) error {
	// Assumes that this is being built from the root directory of the home-automation repo
	args := []string{"build", "-f", dockerfileFilename, "-t", imageName(serviceName), "--rm", "--pull", "."}

	// Setting a pty allows docker to print its progress bars to stdout
	if err := exe.Command("docker", args...).SetPseudoTTY().Run().Err; err != nil {
		return fmt.Errorf("failed to run docker build: %w", err)
	}

	return nil
}

// Stop stops the container for the given service
func Stop(serviceName string) error {
	id, err := containerID(serviceName)
	if err != nil {
		return err
	}

	if id == "" {
		return nil
	}

	return StopByID(id)
}

// Container represents a docker container
type Container struct {
	ID          string
	Image       string
	CreatedAt   string
	RunningFor  string
	Ports       string
	Status      string
	Size        string
	Names       string
	Networks    string
	ServiceName string
}

// GetRunning returns all running containers with the given runScheme label
func GetRunning(runScheme string) ([]*Container, error) {
	args := []string{"ps", "-f", fmt.Sprintf("label=%s=%s", runSchemeLabel, runScheme), "--no-trunc"}

	format := fmt.Sprintf("{{.ID}}§{{.Image}}§{{.CreatedAt}}§{{.RunningFor}}§{{.Ports}}§{{.Status}}§{{.Size}}§{{.Names}}§{{.Networks}}")
	args = append(args, "--format", format)

	result := exe.Command("docker", args...).Run()
	if result.Err != nil {
		return nil, fmt.Errorf("failed to run docker ps: %w", result.Err)
	}

	lines := strings.Split(result.Stdout, "\n")
	var containers []*Container

	for _, line := range lines {
		// Ignore blank lines
		if line == "" {
			continue
		}

		parts := strings.Split(line, "§")

		serviceName, err := serviceNameFromContainerName(parts[7])
		if err != nil {
			return nil, err
		}

		containers = append(containers, &Container{
			ID:          parts[0],
			Image:       parts[1],
			CreatedAt:   parts[2],
			RunningFor:  parts[3],
			Ports:       parts[4],
			Status:      parts[5],
			Size:        parts[6],
			Names:       parts[7],
			Networks:    parts[8],
			ServiceName: serviceName,
		})
	}

	return containers, nil
}

func imageName(serviceName string) string {
	str := strings.ReplaceAll(serviceName, ".", "-")
	str = strings.Replace(str, "service", "s", 1)
	return fmt.Sprintf("%s-%s:%s", "home-automation", str, "latest")
}

// containerName returns the name that should
// be used for a container running this service
func containerName(serviceName string) string {
	// This roughly replicates what docker-compose does
	return fmt.Sprintf("%s_%s", config.Get().ProjectName, serviceName)
}

func serviceNameFromContainerName(containerName string) (string, error) {
	p := fmt.Sprintf("%s_", config.Get().ProjectName)
	if !strings.HasPrefix(containerName, p) {
		return "", fmt.Errorf("container name %q is unexpected", containerName)
	}
	return strings.TrimPrefix(containerName, p), nil
}

// networkName returns the name of the network that containers
// should be a part of. This is the name that docker-compose
// services will join, assuming that no custom network settings
// are specified in the project's docker-compose.yml file.
func networkName() string {
	return fmt.Sprintf("%s_default", config.Get().ProjectName)
}
