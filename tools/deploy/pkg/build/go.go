package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/git"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/deploy/pkg/utils"
)

// GoBuilder is a builder for golang
type GoBuilder struct {
	Service *config.Service
	Target  *config.Target
}

// Build build a go binary for the target architecture and puts it in workingDir
func (b *GoBuilder) Build(revision, workingDir string) (*Release, error) {
	if b.Service.Port == 0 {
		return nil, errors.InternalService("port is not set in config")
	}

	if err := git.Init(revision); err != nil {
		return nil, errors.WithMessage(err, "failed to initialise git mirror")
	}

	op := output.Info("Compiling binary for %s", b.Target.Architecture)

	// Make sure the service exists in the mirror
	pkgToBuild := fmt.Sprintf("./%s/%s", git.MirrorDirectory, b.Service.Name)
	if _, err := os.Stat(filepath.Join(utils.CacheDir(), pkgToBuild)); err != nil {
		op.Failed()
		return nil, errors.WithMessage(err, "failed to stat service directory")
	}

	binName := b.Service.DashedName()

	env := os.Environ()
	switch b.Target.Architecture {
	case config.ArchARMv6:
		env = append(env, "GOOS=linux", "GOARCH=arm", "GOARM=6")
		binName += "-armv6"
	default:
		op.Failed()
		return nil, errors.InternalService("unsupported architecture %q", b.Target.Architecture)
	}

	hash, err := git.CurrentHash(false)
	if err != nil {
		op.Failed()
		return nil, errors.WithMessage(err, "failed to get hash")
	}

	shortHash, err := git.CurrentHash(true)
	if err != nil {
		op.Failed()
		return nil, errors.WithMessage(err, "failed to get short hash")
	}

	binName = fmt.Sprintf("%s-%s", binName, shortHash)
	binOut := filepath.Join(workingDir, binName)

	if err := exe.Command("go", "build", "-o", binOut, pkgToBuild).
		Dir(utils.CacheDir()).Env(env).Run().Err; err != nil {
		op.Failed()
		return nil, errors.WithMessage(err, "failed to compile")
	}

	op.Complete()

	return &Release{
		Cmd:       binName,
		Env:       []*EnvVar{{"PORT", strconv.Itoa(b.Service.Port)}},
		Revision:  hash,
		ShortHash: shortHash,
	}, nil
}
