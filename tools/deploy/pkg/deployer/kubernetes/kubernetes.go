package kubernetes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/build"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/deploy/pkg/utils"
)

const (
	helmReleaseNotFoundErr = "release: not found"
	helmChartPath          = "./tools/deploy/helm/service"
	constPortEnv           = "PORT"
)

// Target is the interface implemented by a Kubernetes target
type Target interface {
	Name() string
	Host() string
	KubeContext() string
	Namespace() string
	DockerRegistry() string
	DockerRepository() string
}

// Kubernetes is a deployer for k8s
type Kubernetes struct {
	Service *config.Service
	Target  Target
}

// Revision returns the currently deployed revision according to the
// revision label on the helm release
func (k *Kubernetes) Revision() (string, error) {
	op := output.Info("Fetching current revision")
	result := exe.Command(
		"helm", "get", "values", k.releaseName(),
		"--namespace", k.Target.Namespace(),
		"--output", "json",
	).Run()
	if result.Err != nil {
		if strings.Contains(result.Err.Error(), helmReleaseNotFoundErr) {
			op.Success()
			return "", nil // Not deployed before
		}

		op.Failed()
		return "", oops.WithMessage(result.Err, "failed to get release values")
	}
	op.Success()

	if len(result.Stderr) > 0 {
		return "", oops.InternalService("helm wrote to stderr: %s", result.Stderr)
	}

	if len(result.Stdout) == 0 {
		return "", oops.InternalService("no response from helm get values")
	}

	values := &struct {
		Revision string `json:"revision"`
	}{}

	if err := json.Unmarshal([]byte(result.Stdout), values); err != nil {
		return "", oops.WithMessage(err, "failed to unmarshal helm response")
	}

	if values.Revision == "" {
		return "", oops.InternalService("could not find revision in helm response: %s", result.Stdout)
	}

	return values.Revision, nil
}

// Deploy upgrades the helm release
func (k *Kubernetes) Deploy(revision string) error {
	builder := build.DockerBuilder{
		Service: k.Service,
		Target:  k.Target,
	}

	dockerBuild, err := builder.Build(revision)
	if err != nil {
		return oops.WithMessage(err, "failed to build docker image")
	}

	if ok, err := k.confirm(dockerBuild); err != nil {
		return oops.WithMessage(err, "failed to confirm")
	} else if !ok {
		return nil
	}

	args := []string{
		"upgrade",
		k.releaseName(),
		helmChartPath,
		"--install",
		"--wait",
		"--kube-context", k.Target.KubeContext(),
		"--namespace", k.Target.Namespace(),
		"--set", "image=" + dockerBuild.Image,
		"--set", "revision=" + dockerBuild.LongHash,
	}

	// If a port is specified in the config, set the containerPort
	// which is used as the targetPort in the k8s service spec.
	if port, ok := dockerBuild.Env.Lookup(constPortEnv); ok {
		args = append(args, "--set", "containerPort="+port)
	}

	for _, v := range dockerBuild.Env {
		s := fmt.Sprintf("config.%s=%s", v.Name, v.Value)
		args = append(args, "--set-string", s)
	}

	op := output.Info("Deploying service")
	if err := exe.Command(
		"helm", args...,
	).Run().Err; err != nil {
		op.Failed()
		return oops.WithMessage(err, "failed to upgrade release")
	}
	op.Success()

	k.success()
	return nil
}

// releaseName returns home-automation-s-foo for service.foo
func (k *Kubernetes) releaseName() string {
	return k.Service.DashedName()
}

func (k *Kubernetes) confirm(db *build.DockerBuild) (bool, error) {
	currentRevision, err := k.Revision()
	if err != nil {
		return false, oops.WithMessage(err, "failed to get current revision")
	}

	return utils.ConfirmDeployment(&utils.Deployment{
		ServiceName:     k.Service.Name(),
		TargetName:      k.Target.Name(),
		TargetHost:      k.Target.Host(),
		CurrentRevision: currentRevision,
		NewRevision:     db.LongHash,
	})
}

func (k *Kubernetes) success() {
	output.InfoLn("\n%s", aurora.Green("Successfully deployed"))
}
