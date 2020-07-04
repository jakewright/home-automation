package env

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/util"
)

// Environment is a set of environment variables
type Environment []*Variable

// Lookup returns the value of the environment variable with the
// given key. If the variable does not exist, the empty string
// and false are returned.
func (e Environment) Lookup(key string) (string, bool) {
	for _, v := range e {
		if v.Name == key {
			return v.Value, true
		}
	}

	return "", false
}

// AsSh returns a slice of the variables in teh Bourne shell format name=value.
func (e Environment) AsSh() []string {
	s := make([]string, len(e))
	for i, v := range e {
		s[i] = v.AsSh()
	}
	return s
}

// Variable represents a single environment variable
type Variable struct {
	Name  string
	Value string
}

// AsSh returns the environment variable in the Bourne shell format name=value.
func (v *Variable) AsSh() string {
	return fmt.Sprintf("%s=%s", v.Name, v.Value)
}

// Parse reads env files into a slice of Variables. If multiple files contain
// the same key, the file specified last takes precedence. Empty lines and
// lines starting with a # are ignored.
func Parse(paths ...string) (Environment, error) {
	m := make(map[string]string)

	for _, p := range paths {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, oops.WithMessage(err, "failed to read %s", p)
		}

		lines := strings.Split(string(b), "\n")
		lines = util.RemoveWhitespaceStrings(lines)

		for _, line := range lines {
			// Skip comments
			if line[0] == '#' {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				return nil, oops.InternalService("unexpected variable %s", line)
			}

			m[parts[0]] = parts[1]
		}
	}

	vars := make([]*Variable, 0, len(m))

	for name, value := range m {
		vars = append(vars, &Variable{
			Name:  name,
			Value: value,
		})
	}

	return vars, nil
}
