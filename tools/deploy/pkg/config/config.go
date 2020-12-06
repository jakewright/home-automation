package config

import (
	"fmt"
)

// Service constants
const (
	LangGo        = "go"
	SysDocker     = "docker"
	SysKubernetes = "kubernetes"
	SysSystemd    = "systemd"
	ArchARMv6     = "ARMv6"
)

var cfg *Config

// Config holds configuration for the deploy tool
type Config struct {
	Repository string
	Targets    map[string]*Target
	Services   map[string]*Service
}

// Get gets the current config, assuming Init() has already been called.
func Get() *Config {
	if cfg == nil {
		panic("config.Init() not called")
	}
	return cfg
}

// Target is the destination server for the deployment
type Target struct {
	// Common
	name   string
	host   string
	system string

	// Systemd
	username     string
	directory    string
	architecture string

	// Kubernetes
	kubeContext      string
	namespace        string
	dockerRegistry   string
	dockerRepository string
}

// Name returns the friendly name of the target
func (t *Target) Name() string {
	return t.name
}

// Host returns the hostname of the target
func (t *Target) Host() string {
	return t.host
}

// System returns the system used by the target
func (t *Target) System() string {
	return t.system
}

// Username returns the username to connect to the target
func (t *Target) Username() string {
	return t.username
}

// Directory returns the working directory to use on the target
func (t *Target) Directory() string {
	return t.directory
}

// Architecture returns the target's architecture
func (t *Target) Architecture() string {
	return t.architecture
}

// KubeContext returns the k8s context to use when interacting with the target
func (t *Target) KubeContext() string {
	return t.kubeContext
}

// Namespace returns the k8s namespace to use when interacting with the target
func (t *Target) Namespace() string {
	return t.namespace
}

// DockerRegistry returns the Docker registry that images should be pushed to
// to be deployed to the target
func (t *Target) DockerRegistry() string {
	return t.dockerRegistry
}

// DockerRepository returns the name of the Docker repository that images should
// be in to be deployed to the target. In the image 192.168.1.1/jakewright/s-foo,
// jakewright is the repository.
func (t *Target) DockerRepository() string {
	return t.dockerRepository
}

// DockerConfig holds options related building the service using Docker
type DockerConfig struct {
	dockerfile string
	args       map[string]string
}

// Dockerfile returns the path to the Dockerfile
func (d *DockerConfig) Dockerfile() string {
	return d.dockerfile
}

// Args returns a map of arguments for the Dockerfile
func (d *DockerConfig) Args() map[string]string {
	return d.args
}

// KubernetesConfig holds options related to a Kubernetes deployment
type KubernetesConfig struct {
	manifests []string
	args      *KubernetesManifestArgs
}

// Manifests returns the paths of the Kubernetes manifest files
func (k *KubernetesConfig) Manifests() []string {
	if k == nil {
		return nil
	}

	return k.manifests
}

// ManifestArgs returns data that should be given to the manifest template
func (k *KubernetesConfig) ManifestArgs() *KubernetesManifestArgs {
	if k == nil {
		return nil
	}

	return k.args
}

// KubernetesManifestArgs holds data that should be given to the manifest template
type KubernetesManifestArgs struct {
	nodePort int
}

// NodePort returns a non-zero int if a custom node port is set in the manifest args
func (a *KubernetesManifestArgs) NodePort() int {
	if a == nil {
		return 0
	}

	return a.nodePort
}

// Service is the microservice to be deployed
type Service struct {
	name        string
	targetNames []string
	targets     []*Target
	language    string
	envFiles    []string
	docker      *DockerConfig
	kubernetes  *KubernetesConfig
}

// Name returns the name of the service
func (s *Service) Name() string {
	return s.name
}

// Path returns the path of the service
func (s *Service) Path() string {
	return "services/" + s.name
}

// TargetNames returns the names of the service's targets
func (s *Service) TargetNames() []string {
	return s.targetNames
}

// Targets returns the service's targets
func (s *Service) Targets() []*Target {
	return s.targets
}

// Language returns the programming language of the service
func (s *Service) Language() string {
	return s.language
}

// EnvFiles returns the env files that should be supplied to the service
// at runtime
func (s *Service) EnvFiles() []string {
	return s.envFiles
}

// Docker returns the docker configuration for the service
func (s *Service) Docker() *DockerConfig {
	return s.docker
}

// Kubernetes returns the Kubernetes configuration for the service
func (s *Service) Kubernetes() *KubernetesConfig {
	return s.kubernetes
}

// DashedName returns home-automation-foo for foo
func (s *Service) DashedName() string {
	// It's important for it to start with home-automation
	// because it's used as the syslog identifier.
	return fmt.Sprintf("home-automation-%s", s.Name())
}

// SyslogIdentifier returns the value that can be used in
// a systemd unit file.
func (s *Service) SyslogIdentifier() string {
	return s.DashedName()
}
