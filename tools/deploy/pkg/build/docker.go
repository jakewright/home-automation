package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/exe"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/config"
	"github.com/jakewright/home-automation/tools/deploy/pkg/git"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
	"github.com/jakewright/home-automation/tools/libraries/env"
)

// reDockerfileArg matches ARG commands in a Dockerfile
var reDockerfileArg = regexp.MustCompile(`ARG ([a-zA-Z0-9-_]+)\n`)

var reDockerImageComponent = regexp.MustCompile(`[a-z0-9]+(?:[._-][a-z0-9]+)*`)

// DockerBuild describes an image that has been built
type DockerBuild struct {
	Image     string
	Env       env.Environment
	LongHash  string
	ShortHash string
}

// DockerDestination is the interface implemented by a Kubernetes target
type DockerDestination interface {
	DockerRegistry() string
	DockerRepository() string
}

// DockerBuilder builds Docker images and pushes them to a registry
type DockerBuilder struct {
	Service *config.Service
	Target  DockerDestination
}

// Build builds the docker image for the given revision and pushes the resulting
// image to the registry defined on the target.
func (b *DockerBuilder) Build(revision string) (*DockerBuild, error) {
	op := output.Info("Preparing build")

	if err := git.Init(revision); err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to initialise git mirror")
	}

	runtimeEnv, err := env.Parse(b.Service.EnvFiles()...)
	if err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to parse service's env files")
	}

	// Make sure the service exists in the mirror
	if _, err := os.Stat(filepath.Join(git.Dir(), b.Service.Path())); err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to stat service directory")
	}

	longHash, shortHash, err := git.CurrentHash()
	if err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to get current hash")
	}

	imageTag := fmt.Sprintf("%s/%s/%s:%s",
		b.Target.DockerRegistry(),
		b.Target.DockerRepository(),
		b.Service.DashedName(),
		shortHash,
	)

	if err := validateDockerImageName(imageTag); err != nil {
		op.Failed()
		return nil, err
	}

	var buildArgs env.Environment
	for name, value := range b.Service.Docker().Args() {
		buildArgs = append(buildArgs, &env.Variable{
			Name:  name,
			Value: value,
		})
	}

	buildArgs = append(buildArgs, &env.Variable{
		Name:  "revision",
		Value: longHash,
	})

	if err := validateBuildArgs(b.Service.Docker(), buildArgs); err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to validate Dockerfile")
	}

	op.Success()

	op = output.Info("Building Docker image")

	dockerArgs := []string{"build"}

	// Append the build args
	for _, buildArg := range buildArgs {
		dockerArgs = append(dockerArgs, "--build-arg", buildArg.AsSh())
	}

	// Append the rest of the arguments
	dockerArgs = append(dockerArgs,
		"-f", b.Service.Docker().Dockerfile(),
		"-t", imageTag,
		"--rm", "--pull", ".")

	cmd := exe.Command("docker", dockerArgs...).
		Dir(git.Dir())

	if output.Verbose {
		cmd = cmd.SetPseudoTTY()
	}

	if err := cmd.Run().Err; err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to build")
	}
	op.Success()

	op = output.Info("Pushing image to registry")
	if err := exe.Command("docker", "push", imageTag).
		Dir(git.Dir()).Run().Err; err != nil {
		op.Failed()
		return nil, oops.WithMessage(err, "failed to push image")
	}
	op.Success()

	return &DockerBuild{
		Image:     imageTag,
		Env:       runtimeEnv,
		LongHash:  longHash,
		ShortHash: shortHash,
	}, nil
}

func validateBuildArgs(conf *config.DockerConfig, args env.Environment) error {
	if conf == nil {
		return oops.InternalService("missing docker config")
	}

	b, err := ioutil.ReadFile(conf.Dockerfile())
	if err != nil {
		return oops.WithMessage(err, "failed to read dockerfile")
	}

	return compareDockerfileArgs(string(b), args)
}

// compareDockerfileArgs returns an error if there is an arg specified in the
// Dockerfile that isn't in the map, or if there is an arg in the map that
// isn't required by the Dockerfile.
func compareDockerfileArgs(dockerfileContent string, givenArgs env.Environment) error {
	matches := reDockerfileArg.FindAllStringSubmatch(dockerfileContent, -1)
	requiredArgs := make(map[string]struct{}, len(matches))

	for _, m := range matches {
		if len(m) != 2 {
			return oops.InternalService("unexpected number of matches %d", len(m))
		}

		requiredArgs[m[1]] = struct{}{}
	}

	for arg := range requiredArgs {
		if _, ok := givenArgs.Lookup(arg); !ok {
			return oops.InternalService("missing Dockerfile arg %q", arg)
		}
	}

	// TODO: support default args that aren't always required, like "revision".
	// for _, arg := range givenArgs {
	// 	if _, ok := requiredArgs[arg.Name]; !ok {
	// 		return oops.InternalService("arg %q specified but not required by Dockerfile", arg)
	// 	}
	// }

	return nil
}

func validateDockerImageName(s string) error {
	if len(s) > 256 {
		return oops.InternalService("docker image name is too long: %s", s)
	}

	for _, part := range strings.Split(s, "/") {
		if !reDockerImageComponent.MatchString(part) {
			return oops.InternalService("invalid docker image name component: %s", part)
		}
	}

	return nil
}
