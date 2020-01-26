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
	parts := strings.Split(path, "/")
	pkg := parts[len(parts)-1]

	return m.addWithPackageName(path, pkg)
}

func (m *importManager) addWithPackageName(path, pkg string) string {
	if path == "" || path == m.pkg {
		return ""
	}

	existing, ok := m.byPath[path]
	if ok {
		return existing.Alias
	}

	if len(m.byPkg[pkg]) > 0 {
		pkg = pkg + strconv.Itoa(len(m.byPkg[pkg]))
	}

	imp := &imp{
		Alias: pkg,
		Path:  path,
	}

	m.byPkg[pkg] = append(m.byPkg[pkg], imp)
	m.byPath[path] = imp

	return pkg
}

func (m *importManager) get() []*imp {
	var imps []*imp
	for _, imp := range m.byPath {
		imps = append(imps, imp)
	}

	return imps
}
