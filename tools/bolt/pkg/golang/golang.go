package golang

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/danielchatfield/go-randutils"

	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/docker"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/toolutils"
)

const runScheme = "golang"

var serviceRe = regexp.MustCompile(`^service\.[a-z0-9-.]+$`)

// System is the golang service management system
type System struct{}

// Is returns true if the given name looks like a service
// name and there's a file at ./[name]/main.go.
func (s *System) Is(name string) (bool, error) {
	// Enforce a pattern on the service name to avoid
	// someone calling build on a random dir that isn't
	// actually a service.
	if !serviceRe.MatchString(name) {
		return false, nil
	}

	fileInfo, err := os.Stat(fmt.Sprintf("./%s/main.go", name))
	if err != nil {
		return false, nil
	}

	// Might as well make sure it's not a directory ¯\_(ツ)_/¯
	return !fileInfo.IsDir(), nil
}

// ListAll returns a list of all service names that _could_ be controlled
// by this system. It doesn't ignore ones that are actually defined
// in the docker-compose.yml file.
func (s *System) ListAll() ([]string, error) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var services []string

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		if is, err := s.Is(f.Name()); err != nil {
			return nil, err
		} else if is {
			services = append(services, f.Name())
		}
	}

	return services, nil
}

// NeedsBuilding returns whether the service needs to be built before it can be run
func (s *System) NeedsBuilding(serviceName string) (bool, error) {
	imageExists, err := docker.ImageForService(serviceName)
	if err != nil {
		return false, fmt.Errorf("failed to check if image for %s exists: %w", serviceName, err)
	}

	return !imageExists, nil
}

// Run starts the service, building first if necessary.
func (s *System) Run(serviceName string) error {
	// Skip if this service is already running
	if running, err := docker.IsRunning(serviceName); err != nil {
		return fmt.Errorf("failed to check status of %s: %w", serviceName, err)
	} else if running {
		return nil
	}

	if needsBuilding, err := s.NeedsBuilding(serviceName); err != nil {
		return err
	} else if needsBuilding {
		if err := s.Build(serviceName); err != nil {
			return fmt.Errorf("failed to build %s: %w", serviceName, err)
		}
	}

	if err := docker.Run(serviceName, runScheme); err != nil {
		return fmt.Errorf("failed to run %s container: %w", serviceName, err)
	}

	return nil
}

// IsRunning returns whether the service is running
func (s *System) IsRunning(serviceName string) (bool, error) {
	return docker.IsRunning(serviceName)
}

// Build builds the service
func (s *System) Build(serviceName string) error {
	tmpl, err := template.ParseFiles(config.Get().GoDockerfileTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile template: %w", err)
	}

	// Create a Dockerfile from the template
	data := &struct {
		GoVersion, Service string
	}{config.Get().GoVersion, serviceName}

	b := bytes.Buffer{}
	if err := tmpl.Execute(&b, data); err != nil {
		return fmt.Errorf("failed to execute Dockerfile template: %w", err)
	}

	// Create a filename. Include random characters so this doesn't
	// conflict with any simultaneous builds for the same service.
	rand, err := randutils.String(6)
	if err != nil {
		return fmt.Errorf("failed to generate random string: %w", err)
	}
	dockerfileFilename := filepath.Join(toolutils.CacheDir(), fmt.Sprintf("%s-%s.dockerfile", serviceName, strings.ToLower(rand)))

	// Write the Dockerfile to the cache directory
	if err := ioutil.WriteFile(dockerfileFilename, b.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}

	if err := docker.Build(serviceName, dockerfileFilename); err != nil {
		return fmt.Errorf("failed to run docker build: %w", err)
	}

	if err := os.Remove(dockerfileFilename); err != nil {
		return fmt.Errorf("failed to remove Dockerfile: %w", err)
	}

	return nil
}

// Stop stops a service
func (s *System) Stop(serviceName string) error {
	return docker.Stop(serviceName)
}

// StopAll stops all golang services
func (s *System) StopAll() error {
	containers, err := docker.GetRunning(runScheme)
	if err != nil {
		return err
	}

	for _, c := range containers {
		op := output.Info("Stopping %s", c.ServiceName)
		// Stopping by ID instead of service name is an optimisation
		if err := docker.StopByID(c.ID); err != nil {
			op.Failed()
			return err
		}
		op.Complete()
	}

	return nil
}

// Exec isn't needed for golang services
func (s *System) Exec(_, _ string, _ string, _ ...string) error {
	panic("not implemented")
}
