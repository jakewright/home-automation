package imports

import (
	"strconv"
	"strings"
)

type Imp struct {
	Alias string
	Path  string
}

type Manager struct {
	// self is the package for which we're managing
	// imports. It should be a full import path e.g.
	// github.com/jakewright/home-automation/service.foo
	self   string
	byPath map[string]*Imp
	byPkg  map[string][]*Imp
}

func NewManager(self string) *Manager {
	return &Manager{
		self:   self,
		byPath: make(map[string]*Imp),
		byPkg:  make(map[string][]*Imp),
	}
}

func (m *Manager) Add(path string) string {
	parts := strings.Split(path, "/")
	pkg := parts[len(parts)-1]

	return m.addWithPackageName(path, pkg)
}

func (m *Manager) addWithPackageName(path, pkg string) string {
	if path == "" || path == m.self {
		return ""
	}

	existing, ok := m.byPath[path]
	if ok {
		return existing.Alias
	}

	if len(m.byPkg[pkg]) > 0 {
		pkg = pkg + strconv.Itoa(len(m.byPkg[pkg]))
	}

	imp := &Imp{
		Alias: pkg,
		Path:  path,
	}

	m.byPkg[pkg] = append(m.byPkg[pkg], imp)
	m.byPath[path] = imp

	return pkg
}

func (m *Manager) Get() []*Imp {
	var imps []*Imp
	for _, imp := range m.byPath {
		imps = append(imps, imp)
	}

	return imps
}
