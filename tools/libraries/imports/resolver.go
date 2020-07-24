package imports

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// Resolver figures out go import paths for packages
type Resolver struct {
	// module is the go module name e.g. github.com/jakewright/home-automation
	module string

	// moduleRoot is the absolute path to the root of the go module
	moduleRoot string
}

// NewResolver figures out the go module name
// and directory and returns an initialised Resolver
func NewResolver() (*Resolver, error) {
	// In order to figure out import paths, the enricher needs
	// to know the go module name and the path to the module root.
	var module string
	var modFilePath string

	// Iteratively look for a go.mod file in subsequently higher parent directories.
	// The maximum depth (height?) of 10 is arbitrary and could be increased
	// if this ends up being used in a very nested directory structure.
	for i := 0; i < 10; i++ {
		modFilePath = strings.Repeat("../", i) + "go.mod"
		if i == 0 {
			modFilePath = "./" + modFilePath
		}

		// Try to read the file to see if it exists
		b, err := ioutil.ReadFile(modFilePath)
		if err != nil {
			// If it didn't exist, try again
			if os.IsNotExist(err) {
				continue
			}
			return nil, err // Unexpected error
		}

		// Pull out the module name e.g. github.com/jakewright/home-automation
		module = modfile.ModulePath(b)
		break
	}

	if module == "" {
		return nil, fmt.Errorf("failed to find module path")
	}

	// Get an absolute path to the root of the go module
	moduleRoot, err := filepath.Abs(filepath.Dir(modFilePath))
	if err != nil {
		return nil, err
	}

	return &Resolver{
		module:     module,
		moduleRoot: moduleRoot,
	}, nil
}

// Resolve returns the import path for a package in service s.
// It assumes that the def file exists in the root of service s.
func (r *Resolver) Resolve(defPath, pkg string) (string, error) {
	// Make sure the .def path is an absolute path
	var err error
	defPath, err = filepath.Abs(defPath)
	if err != nil {
		return "", err
	}

	defPathRelToRoot, err := filepath.Rel(r.moduleRoot, defPath)
	if err != nil {
		return "", err
	}

	// Prepend the module name and remove the /xxx.def from the end
	importPath := filepath.Dir(filepath.Join(r.module, defPathRelToRoot))

	// Append the package path
	return filepath.Join(importPath, pkg), nil
}
