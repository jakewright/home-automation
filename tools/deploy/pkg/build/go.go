package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/git"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/libraries/env"
)

const revisionVarPath = "github.com/jakewright/home-automation/libraries/go/bootstrap.Revision"

// GoBuilder is a builder for golang
type GoBuilder struct {
	Service *config.Service
	Target  Machine
}

var _ LocalBuilder = (*GoBuilder)(nil)

// Build build a go binary for the target architecture and puts it in workingDir
func (b *GoBuilder) Build(revision, workingDir string) (*Release, error) {
	if err := git.Init(revision); err != nil {
		return nil, oops.WithMessage(err, "failed to initialise git mirror")
	}

	op := output.Info("Parsing service's config")
	runtimeEnv, err := env.Parse(b.Service.EnvFiles()...)
	if err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to parse service's env files")
	}
	op.Success()

	op = output.Info("Compiling binary for %s", b.Target.Architecture())

	// Make sure the service exists in the mirror
	pkgToBuild := b.Service.Path()
	if _, err := os.Stat(filepath.Join(git.Dir(), pkgToBuild)); err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to stat service directory")
	}

	binName := b.Service.DashedName()

	buildEnv := os.Environ()
	switch b.Target.Architecture() {
	case config.ArchARMv6:
		buildEnv = append(buildEnv, "GOOS=linux", "GOARCH=arm", "GOARM=6")
		binName += "-armv6"
	default:
		op.Failed()
		return nil, oops.InternalService("unsupported architecture %q", b.Target.Architecture())
	}

	hash, shortHash, err := git.CurrentHash()
	if err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to get current hash")
	}

	binName = fmt.Sprintf("%s-%s", binName, shortHash)
	binOut := filepath.Join(workingDir, binName)

	flags := fmt.Sprintf("-X %s=%s", revisionVarPath, hash)

	if err := exe.Command("go", "build", "-o", binOut, "-ldflags", flags, pkgToBuild).
		Dir(git.Dir()).Env(buildEnv).Run().Err; err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to compile")
	}

	op.Success()

	return &Release{
		Cmd:       binName,
		Env:       runtimeEnv,
		Revision:  hash,
		ShortHash: shortHash,
	}, nil
}
