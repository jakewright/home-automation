package kubernetes

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/logrusorgru/aurora"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/build"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/deploy/pkg/utils"
)

const (
	k8sNotFoundErr = "NotFound"
	constPortEnv   = "PORT"
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

type manifestTemplateData struct {
	// ServiceName becomes the name of the Kubernetes
	// service and deployment
	ServiceName string

	// Image is the name of the Docker image the pods should run
	Image string

	// Revision should be the full git hash
	Revision string

	// DeploymentTimestamp should be set to time.Now() when
	// the deployment happens. This makes sure that each
	// deployment contains a change, even if the image name
	// hasn't changed. Combined with imagePullPolicy: Always,
	// all deployments cause K8s to re-pull the image and
	// roll the pods.
	DeploymentTimestamp string

	// ServicePort is the port that will be exposed by the service
	ServicePort int

	// ContainerPort is the port that the service should
	// access in the pods that it targets
	ContainerPort int

	NodePort int

	// Config is a map of environment variables that the
	// service should have at runtime
	Config map[string]string
}

// Kubernetes is a deployer for k8s
type Kubernetes struct {
	Service *config.Service
	Target  Target
}

// Revision returns the currently deployed revision from the k8s annotation
func (k *Kubernetes) Revision() (string, error) {
	op := output.Info("Fetching current revision")
	result := exe.Command(
		"kubectl",
		"get",
		fmt.Sprintf("deployments/%s", k.releaseName()),
		"--context", k.Target.KubeContext(),
		"--namespace", k.Target.Namespace(),
		"--output", "jsonpath=\"{.metadata.annotations.revision}\"",
	).Run()
	if result.Err != nil {
		if strings.Contains(result.Err.Error(), k8sNotFoundErr) {
			op.Success()
			return "", nil // Not deployed before
		}

		op.Failed()
		return "", oops.WithMessage(result.Err, "failed to get release values")
	}
	op.Success()

	if len(result.Stderr) > 0 {
		return "", oops.InternalService("kubectl wrote to stderr: %s", result.Stderr)
	}

	if len(result.Stdout) == 0 {
		return "", oops.InternalService("no response from kubectl")
	}

	revision, err := strconv.Unquote(result.Stdout)
	if err != nil {
		return "", oops.WithMessage(err, "failed to read response")
	}

	return revision, nil
}

// Deploy builds and applies a kubernetes manifest for the service
func (k *Kubernetes) Deploy(revision string) error {
	if len(k.Service.Kubernetes().Manifests()) == 0 {
		return oops.InternalService("no k8s manifests specified for service")
	}

	builder := build.DockerBuilder{
		Service: k.Service,
		Target:  k.Target,
	}

	dockerBuild, err := builder.Build(revision)
	if err != nil {
		return oops.WithMessage(err, "failed to build docker image")
	}

	args := &manifestTemplateData{
		ServiceName:         k.releaseName(),
		Image:               dockerBuild.Image,
		Revision:            dockerBuild.LongHash,
		DeploymentTimestamp: time.Now().Format(time.RFC3339),
		ServicePort:         80,
		ContainerPort:       80, // Might be overridden below
		NodePort:            k.Service.Kubernetes().ManifestArgs().NodePort(),
		Config:              make(map[string]string),
	}

	// If a port is specified in the config, set the containerPort
	// which is used as the targetPort in the k8s service spec.
	if port, ok := dockerBuild.Env.Lookup(constPortEnv); ok {
		p, err := strconv.Atoi(port)
		if err != nil {
			return oops.WithMessage(err, "failed to parse port from service's env file")
		}
		args.ContainerPort = p
	}

	for _, v := range dockerBuild.Env {
		args.Config[v.Name] = v.Value
	}

	var manifest bytes.Buffer

	// Read the k8s resource files
	for _, filename := range k.Service.Kubernetes().Manifests() {
		t, err := template.ParseFiles(filename)
		if err != nil {
			return oops.WithMessage(err, "failed to parse manifest: %q", filename)
		}

		manifest.WriteString("---\n")

		if err := t.Execute(&manifest, args); err != nil {
			return oops.WithMessage(err, "failed to execute template")
		}
	}

	op := output.Info("Deploying service")

	// TODO: use the go library instead of kubectl
	cmd := exe.Command(
		"kubectl",
		"apply",
		"--filename", "-", // Read from stdin
		"--context", k.Target.KubeContext(),
		"--namespace", k.Target.Namespace(),
	).SetInput(manifest.String())

	output.Debug("\n\n%s\n\n", manifest.String())

	if ok, err := k.confirm(dockerBuild); err != nil {
		return oops.WithMessage(err, "failed to confirm")
	} else if !ok {
		return nil
	}

	res := cmd.Run()

	if res.Err != nil {
		op.Failed()
		return oops.WithMessage(res.Err, "failed to upgrade release")
	}

	output.Info(res.Stdout)

	op.Success()

	k.success()
	return nil
}

// releaseName returns home-automation-foo for foo
func (k *Kubernetes) releaseName() string {
	return k.Service.Name()
}

func (k *Kubernetes) confirm(db *build.DockerBuild) (bool, error) {
	currentRevision, err := k.Revision()
	if err != nil {
		return false, oops.WithMessage(err, "failed to get current revision")
	}

	return utils.ConfirmDeployment(&utils.Deployment{
		ServiceName:     k.Service.Name(),
		ServicePath:     k.Service.Path(),
		TargetName:      k.Target.Name(),
		TargetHost:      k.Target.Host(),
		CurrentRevision: currentRevision,
		NewRevision:     db.LongHash,
	})
}

func (k *Kubernetes) success() {
	output.InfoLn("\n%s", aurora.Green("Successfully deployed"))
}
