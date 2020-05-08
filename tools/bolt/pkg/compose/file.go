package compose

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type composeFile struct {
	Version  string                     `yaml:"version"`
	Services map[string]*composeService `yaml:"services"`
	Networks map[string]interface{}     `yaml:"networks"`
}

type composeService struct {
	Image string   `yaml:"image"`
	Ports []string `yaml:"ports"`
}

func parse(filename string) (*composeFile, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read docker-compose file: %w", err)
	}

	f := &composeFile{}
	if err := yaml.Unmarshal(b, f); err != nil {
		return nil, fmt.Errorf("failed to unmarshal docker-compose composeFile: %w", err)
	}

	return f, nil
}
