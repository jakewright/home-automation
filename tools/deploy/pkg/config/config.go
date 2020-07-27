package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

// Service constants
const (
	LangGo         = "go"
	LangJavaScript = "javascript"
	SysDocker      = "docker"
	SysKubernetes  = "kubernetes"
	SysSystemd     = "systemd"
	ArchARMv6      = "ARMv6"
)

var cfg config

type config struct {
	Repository string              `yaml:"repository"`
	Targets    map[string]*Target  `yaml:"targets"`
	Services   map[string]*Service `yaml:"services"`
}

// Target is the destination server for the deployment
type Target struct {
	// Common
	Name   string `yaml:"-"`
	Host   string `yaml:"host"`
	System string `yaml:"system"`

	// Systemd
	Username     string `yaml:"username"`
	Directory    string `yaml:"directory"`
	Architecture string `yaml:"architecture"`

	// Kubernetes
	KubeContext      string `yaml:"kube_context"`
	Namespace        string `yaml:"namespace"`
	DockerRegistry   string `yaml:"docker_registry"`
	DockerRepository string `yaml:"docker_repository"`
}

// DockerConfig holds options related building the service using Docker
type DockerConfig struct {
	Dockerfile string            `yaml:"dockerfile"`
	Args       map[string]string `yaml:"args"`
}

// Service is the microservice to be deployed
type Service struct {
	Name        string        `yaml:"-"`
	TargetNames []string      `yaml:"targets"`
	Targets     []*Target     `yaml:"-"`
	Language    string        `yaml:"language"`
	EnvFiles    []string      `yaml:"env_files"`
	Docker      *DockerConfig `yaml:"docker"`
}

// DashedName returns home-automation-s-foo for service.foo
func (s *Service) DashedName() string {
	str := strings.ReplaceAll(s.Name, ".", "-")
	str = strings.Replace(str, "service", "s", 1)
	// It's important for it to start with home-automation
	// because it's used as the syslog identifier.
	return fmt.Sprintf("home-automation-%s", str)
}

// SyslogIdentifier returns the value that can be used in
// a systemd unit file.
func (s *Service) SyslogIdentifier() string {
	return s.DashedName()
}

// Init reads and validates config
func Init(filename string) (err error) {
	op := output.Info("Reading config from %v", filename)
	defer func() {
		if err == nil {
			op.Success()
		} else {
			op.Failed()
		}
	}()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return oops.WithMessage(err, "failed to read config file")
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return oops.WithMessage(err, "failed to unmarshal config")
	}

	for name, target := range cfg.Targets {
		target.Name = name

		switch target.System {
		case SysDocker, SysKubernetes, SysSystemd: // ok
		default:
			return oops.InternalService("Invalid system %q for target %q", target.System, name)
		}

		switch target.Architecture {
		case "", ArchARMv6: // ok
		default:
			return oops.InternalService("Invalid architecture %q for target %q", target.Architecture, name)
		}
	}
	for name, service := range cfg.Services {
		service.Name = name

		switch service.Language {
		case LangGo, LangJavaScript: // ok
		default:
			return oops.InternalService("Invalid language '%s' for service '%s'", service.Language, name)
		}

		for _, targetName := range service.TargetNames {
			target := findTarget(targetName)
			if target == nil {
				return oops.InternalService("Invalid target %q for service %q", targetName, name)
			}

			service.Targets = append(service.Targets, target)
		}
	}

	return nil
}

func findTarget(name string) *Target {
	return cfg.Targets[name]
}

// FindService returns a service by name or nil if it doesn't exist
func FindService(name string) *Service {
	return cfg.Services[name]
}

// Repository returns the repo specified in config
func Repository() string {
	return cfg.Repository
}
