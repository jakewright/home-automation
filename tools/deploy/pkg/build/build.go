package build

import (
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/libraries/env"
)

// Release represents something that can be deployed
type Release struct {
	Cmd       string
	Env       env.Environment
	Revision  string
	ShortHash string
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

	return nil, oops.BadRequest("no suitable builder")
}
