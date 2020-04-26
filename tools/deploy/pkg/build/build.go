package build

import (
	"fmt"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
)

// Release represents something that can be deployed
type Release struct {
	Cmd       string
	Env       []*EnvVar
	ShortHash string
}

// EnvVar represents a single environment variable
type EnvVar struct {
	Name  string
	Value string
}

// AsSh returns the environment variable in the Bourne shell format name=value.
func (e *EnvVar) AsSh() string {
	return fmt.Sprintf("%s=%s", e.Name, e.Value)
}

// Builder prepares a release
type Builder interface {
	Build(revision, workingDir string) (*Release, error)
}

// Choose returns a builder based on the service and target
func Choose(service *config.Service, target *config.Target) (Builder, error) {
	switch {
	case service.Language == config.LangGo && target.System == config.SysSystemd:
		return &GoBuilder{
			Service: service,
			Target:  target,
		}, nil
	}

	return nil, errors.BadRequest("no suitable builder")
}
