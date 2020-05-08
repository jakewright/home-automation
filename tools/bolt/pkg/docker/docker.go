package docker

import (
	"fmt"
	"strconv"

	"github.com/jakewright/home-automation/libraries/go/exe"
)

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
