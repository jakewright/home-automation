package docker

import (
	"fmt"
	"strconv"

	"github.com/jakewright/home-automation/libraries/go/exe"
)

/* These functions do not provide the service-level abstraction that the other
files provide, i.e. they operate directly on image and container names. */

// IsContainerRunning returns whether the container with
// the given ID is currently in a running state.
func IsContainerRunning(id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("id is empty")
	}

	result := exe.Command("docker", "inspect", "-f", "{{.State.Running}}", id).Run()
	if result.Err != nil {
		return false, fmt.Errorf("failed to run docker inspect: %w", result.Err)
	}

	b, err := strconv.ParseBool(result.Stdout)
	if err != nil {
		return false, fmt.Errorf("failed to parse docker inspect output: %w", err)
	}

	return b, nil
}

// StopByID stops the container with the given ID
func StopByID(id string) error {
	if err := exe.Command("docker", "stop", id).Run().Err; err != nil {
		return fmt.Errorf("failed to run docker stop: %w", err)
	}

	return nil
}

// ImageExists returns whether the image with the given name exists
func ImageExists(imageName string) (bool, error) {
	result := exe.Command("docker", "images", "-q", imageName).Run()
	if result.Err != nil {
		return false, fmt.Errorf("failed to run docker images: %w", result.Err)
	}

	return len(result.Stdout) > 0, nil
}
