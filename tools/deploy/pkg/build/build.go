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

// LocalBuilder prepares a release
type LocalBuilder interface {
	Build(revision, workingDir string) (*Release, error)
}

// ChooseLocal returns a builder based on the service and target
func ChooseLocal(service *config.Service, target *config.Target) (LocalBuilder, error) {
	switch service.Language {
	case config.LangGo:
		return &GoBuilder{
			Service: service,
			Target:  target,
		}, nil
	}

	return nil, oops.BadRequest("no suitable builder")
}
