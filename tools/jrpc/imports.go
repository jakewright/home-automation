package main

import (
	"strconv"
	"strings"
)

type importManager struct {
	// pkg is the package for which we're managing
	// imports. It should be a full import path e.g.
	// github.com/jakewright/home-automation/service.foo
	pkg    string
	byPath map[string]*imp
	byPkg  map[string][]*imp
}

func newImportManager(pkg string) *importManager {
	return &importManager{
		pkg:    pkg,
		byPath: make(map[string]*imp),
		byPkg:  make(map[string][]*imp),
	}
}

func (m *importManager) add(path string) string {
	if path == "" || path == m.pkg {
		return ""
	}

	parts := strings.Split(path, "/")
	pkg := parts[len(parts)-1]

	existing, ok := m.byPath[path]
	if ok {
		if existing.Alias != "" {
			return existing.Alias
		}

		return pkg
	}

	var alias string
	if len(m.byPkg[pkg]) > 0 {
		alias = pkg + strconv.Itoa(len(m.byPkg[pkg]))
	}

	imp := &imp{
		Alias: alias,
		Path:  path,
	}

	m.byPkg[pkg] = append(m.byPkg[pkg], imp)
	m.byPath[path] = imp

	if alias != "" {
		return alias
	}

	return pkg
}

func (m *importManager) get() []*imp {
	var imps []*imp
	for _, imp := range m.byPath {
		imps = append(imps, imp)
	}

	return imps
}
