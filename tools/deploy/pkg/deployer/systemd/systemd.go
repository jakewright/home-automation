package systemd

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/build"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/deploy/pkg/utils"
)

// Target is the interface implemented by a systemd target
type Target interface {
	Name() string
	Host() string
	Username() string
	Directory() string
	Architecture() string
}

// Systemd is a deployer for services running under systemd
type Systemd struct {
	Service *config.Service
	Target  Target
}

// Revision returns the currently deployed revision
func (d *Systemd) Revision() (string, error) {
	op := output.Info("Connecting to %s", d.Target.Host())
	client, err := exe.SSHClient(d.Target.Username(), d.Target.Host())
	if err != nil {
		op.Failed()
		return "", oops.WithMessage(err, "failed to create SSH client")
	}
	defer func() { _ = client.Close() }()
	op.Success()

	cmd := fmt.Sprintf("sudo systemctl cat %s.service", d.Service.DashedName())

	result := exe.RemoteCommand(cmd).Run(client)

	if strings.HasPrefix(result.Stderr, "No files found for") {
		return "", nil // Not deployed before
	}

	if result.Err != nil {
		return "", oops.WithMessage(err, "failed to cat service file")
	}

	re := regexp.MustCompile(`X-Revision=([a-zA-z0-9]+)\n`)
	matches := re.FindStringSubmatch(result.Stdout)

	if len(matches) != 2 {
		return "", oops.InternalService("unexpected match length %d", len(matches))
	}

	return matches[1], nil
}

// Deploy deploys the given revision of the service to the target
func (d *Systemd) Deploy(revision string) error {
	if !filepath.IsAbs(d.Target.Directory()) {
		return oops.InternalService("target directory is not absolute: %s", d.Target.Directory)
	}

	deploymentName, err := d.generateDeploymentName(revision)
	if err != nil {
		return oops.WithMessage(err, "failed to generate deployment name")
	}

	workingDir, err := d.workingDir(deploymentName)
	if err != nil {
		return oops.WithMessage(err, "failed to create working directory")
	}

	builder, err := build.ChooseLocal(d.Service, d.Target)
	if err != nil {
		return oops.WithMessage(err, "failed to choose builder")
	}

	release, err := builder.Build(revision, workingDir)
	if err != nil {
		return oops.WithMessage(err, "failed to build")
	}

	if ok, err := d.confirm(release); err != nil {
		return oops.WithMessage(err, "failed to confirm")
	} else if !ok {
		return nil
	}

	op := output.Info("Connecting to %s", d.Target.Host())
	client, err := exe.SSHClient(d.Target.Username(), d.Target.Host())
	if err != nil {
		op.Failed()
		return oops.WithMessage(err, "failed to create SSH client")
	}
	defer func() { _ = client.Close() }()
	op.Success()

	op = output.Info("Copying files to %s", d.Target.Host())
	if err := utils.SCP(workingDir, d.Target.Username(), d.Target.Host(), d.Target.Directory()); err != nil {
		op.Failed()
		return oops.WithMessage(err, "failed to copy binary to target")
	}
	op.Success()

	op = output.Info("Updating unit file")

	if err := d.updateUnitFile(client, release, deploymentName); err != nil {
		op.Failed()
		return oops.WithMessage(err, "failed to update unit file")
	}
	op.Success()

	op = output.Info("Restarting service")
	if err := d.restartUnit(client); err != nil {
		op.Failed()
		return oops.WithMessage(err, "failed to enable service")
	}
	op.Success()

	op = output.Info("Cleaning up old files")
	if err := d.cleanup(client, deploymentName, workingDir); err != nil {
		op.Failed()
		return oops.WithMessage(err, "failed to clean up old files")
	}
	op.Success()

	d.success(release)

	return nil
}
