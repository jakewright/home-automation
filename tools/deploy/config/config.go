package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const (
	LangGo = "go"
	LangJavaScript = "javascript"
	SysDocker = "docker"
	SysSystemd = "systemd"
)

var cfg config

type config struct {
	Targets map[string]*Target `yaml:"targets"`
	Services map[string]*Service `yaml:"services"`
}

type Target struct {
	Name string `yaml:"-"`
	Host string `yaml:"host"`
	Username string `yaml:"username"`
	Directory string `yaml:"directory"`
}

type Service struct {
	Name string `yaml:"-"`
	Target string `yaml:"target"`
	Language string `yaml:"language"`
	System string `yaml:"system"`
}

func init() {
	b, err := ioutil.ReadFile("./private/deploy/config.yml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v\n", err)
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v\n", err)
	}

	for name, target := range cfg.Targets {
		target.Name = name
	}
	for name, service := range cfg.Services {
		service.Name = name

		switch service.Language {
		case LangGo, LangJavaScript: // ok
		default:
			log.Fatalf("Invalid language '%s' for service '%s'\n", service.Language, name)
		}

		switch service.System {
		case SysDocker, SysSystemd: // ok
		default:
			log.Fatalf("Invalid system '%s' for service '%s'\n", service.System, name)
		}
	}
}

func FindTarget(name string) *Target {
	return cfg.Targets[name]
}

func FindService(name string) *Service {
	return cfg.Services[name]
}