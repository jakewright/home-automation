package systemd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/danielchatfield/go-randutils"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/build"
	"github.com/jakewright/home-automation/tools/deploy/pkg/git"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/deploy/pkg/utils"
	"github.com/jakewright/home-automation/tools/libraries/cache"
)

// generateDeploymentName returns a unique name for this deployment. It ends in
// random characters to allow the same revision to be deployed to the same target
// multiple times.
func (d *Systemd) generateDeploymentName(revision string) (string, error) {
	if err := git.Init(revision); err != nil {
		return "", oops.WithMessage(err, "failed to initialise git mirror")
	}

	_, shortHash, err := git.CurrentHash()
	if err != nil {
		return "", oops.WithMessage(err, "failed to get short hash")
	}

	random, err := randutils.String(4)
	if err != nil {
		return "", oops.WithMessage(err, "failed to generate random string")
	}

	return fmt.Sprintf("%s-%s-%s-%s", d.Service.DashedName(), d.Target.Name(), shortHash, random), nil
}

// generateDeploymentGlob returns a glob that can be passed to rm
// on the target to remove all versions of the service
func (d *Systemd) generateDeploymentGlob() string {
	// This should be the prefix of the working directory because the whole
	// thing gets copied to the target and keeps the same name.
	return fmt.Sprintf("%s-%s*", d.Service.DashedName(), d.Target.Name())
}

// workingDir creates and returns the full path to a temporary
// working directory on the local filesystem.
func (d *Systemd) workingDir(deploymentName string) (string, error) {
	workingDir := filepath.Join(cache.Dir(), "deployments", deploymentName)

	if err := os.MkdirAll(workingDir, os.ModePerm); err != nil {
		return "", oops.WithMessage(err, "failed to create working directory")
	}

	return workingDir, nil
}

func (d *Systemd) confirm(release *build.Release) (bool, error) {
	currentRevision, err := d.Revision()
	if err != nil {
		return false, oops.WithMessage(err, "failed to get current revision")
	}

	return utils.ConfirmDeployment(&utils.Deployment{
		ServiceName:     d.Service.Name(),
		ServicePath:     d.Service.Path(),
		TargetName:      d.Target.Name(),
		TargetHost:      d.Target.Host(),
		CurrentRevision: currentRevision,
		NewRevision:     release.Revision,
	})
}

func (d *Systemd) updateUnitFile(client *ssh.Client, release *build.Release, deploymentName string) error {
	unit, err := d.unitFile(deploymentName, release)
	if err != nil {
		return oops.WithMessage(err, "failed to write unit file")
	}

	// The -e flag on echo interprets backslash escapes. This is piped into
	// systemctl edit, which by default opens an editor but this can be
	// overridden with SYSTEMD_EDITOR. Using tee will write stdin to the file.
	// The --full flag overwrites the entire unit file and --force will create
	// a unit file if it didn't previously exist, so this works for new
	// services and for services that have been previously deployed. Systemctl
	// reloads its config afterwards in a way that is equivalent to daemon-reload.

	if err := exe.RemoteCommand(
		fmt.Sprintf(
			"echo -e %q | sudo SYSTEMD_EDITOR=tee systemctl edit --full --force %s",
			unit,
			d.Service.DashedName(),
		)).RequestPseudoTTY().Run(client).Err; err != nil {
		return oops.WithMessage(err, "failed to edit unit file")
	}

	return nil
}

func (d *Systemd) unitFile(deploymentName string, release *build.Release) ([]byte, error) {
	txt := `[Unit]
Description={{ .Description }}

[Service]
SyslogIdentifier={{ .SyslogIdentifier }}
WorkingDirectory={{ .WorkingDirectory }}
{{- range .Environment }}
Environment={{ . }}
{{- end }}
Type=idle
ExecStart={{ .ExecStart }}
X-Revision={{ .Revision }}

[Install]
WantedBy=multi-user.target
`

	tmpl, err := template.New("unit").Parse(txt)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to parse template")
	}

	cmd := filepath.Join(d.Target.Directory(), deploymentName, release.Cmd)
	env := make([]string, len(release.Env))
	for i, e := range release.Env {
		env[i] = e.AsSh()
	}

	data := struct {
		Description, SyslogIdentifier, WorkingDirectory, ExecStart, Revision string
		Environment                                                          []string
	}{
		Description:      d.Service.Name(),
		SyslogIdentifier: d.Service.SyslogIdentifier(),
		WorkingDirectory: d.Target.Directory(),
		Environment:      env,
		ExecStart:        cmd,
		Revision:         release.Revision,
	}

	b := bytes.Buffer{}
	if err := tmpl.Execute(&b, data); err != nil {
		return nil, oops.WithMessage(err, "failed to execute template")
	}

	return b.Bytes(), nil
}

func (d *Systemd) restartUnit(client *ssh.Client) error {
	if err := exe.RemoteCommand("sudo systemctl enable", d.Service.DashedName()).Run(client).Err; err != nil {
		return oops.WithMessage(err, "failed to enable service")
	}

	if err := exe.RemoteCommand("sudo systemctl restart", d.Service.DashedName()).Run(client).Err; err != nil {
		return oops.WithMessage(err, "failed to restart service")
	}

	return nil
}

func (d *Systemd) cleanup(client *ssh.Client, deploymentName, workingDir string) error {
	// Delete all old deployment folders from the target
	if err := exe.RemoteCommand(fmt.Sprintf(
		"find %s -maxdepth 1 -name '%s' ! -name '%s' -type d -exec rm -r {} +",
		d.Target.Directory(),
		d.generateDeploymentGlob(),
		deploymentName,
	)).Run(client).Err; err != nil {
		return err
	}

	// Delete the working directory
	if err := exe.Command("rm", "-r", workingDir).Run().Err; err != nil {
		return err
	}

	return nil
}

func (d *Systemd) success(release *build.Release) {
	output.InfoLn("\n%s", aurora.Green("Successfully deployed"))
	port, _ := release.Env.Lookup("PORT")
	svc := aurora.Sprintf(aurora.Index(105, "http://%s:%s/"), d.Target.Host, port)
	output.InfoLn("Service available at %s\n", svc)
}
