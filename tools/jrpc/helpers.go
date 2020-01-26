package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/jakewright/home-automation/libraries/go/svcdef"
)

func messagesFromFile(opts *options, file *svcdef.File, im *importManager) ([]*message, error) {
	var messages []*message

	for _, m := range file.Messages {
		message, err := newMessage(m, im, opts, file)
		if err != nil {
			return nil, err
		}

		for _, nm := range m.Nested {
			nested, err := newMessage(nm, im, opts, file)
			if err != nil {
				return nil, err
			}

			messages = append(messages, nested)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func newMessage(m *svcdef.Message, im *importManager, opts *options, file *svcdef.File) (*message, error) {
	name := m.QualifiedName
	name = strings.ReplaceAll(name, ".", "_")

	re := regexp.MustCompile("^[A-Z][a-zA-Z0-9_]*$")
	if !re.MatchString(name) {
		return nil, fmt.Errorf("invalid message name %s", name)
	}

	fields := make([]*field, len(m.Fields))
	for j, f := range m.Fields {
		goName, jsonName, err := getGoFieldName(f.Name)
		if err != nil {
			return nil, err
		}

		goTypeName, err := fieldToType(f, opts.DefPath, file.Imports, im)
		if err != nil {
			return nil, fmt.Errorf("failed to get field type in message %s: %v", m.QualifiedName, err)
		}

		isMessageType := strings.HasPrefix(f.QualifiedType, ".")
		pointer := strings.HasPrefix(goTypeName, "*")

		var required bool
		if v, ok := f.Options["required"].(bool); ok {
			required = v
		}

		fields[j] = &field{
			GoName:        goName,
			JSONName:      jsonName,
			Type:          goTypeName,
			IsMessageType: isMessageType,
			Repeated:      f.Repeated,
			Required:      required,
			Pointer:       pointer,
		}
	}

	return &message{
		Name:   name,
		Fields: fields,
	}, nil
}

func fieldToType(field *svcdef.Field, defPath string, imports map[string]*svcdef.Import, im *importManager) (string, error) {
	typ, err := resolveTypeName(field.QualifiedType, defPath, imports, im)
	if err != nil {
		return "", fmt.Errorf("failed to resolve type in field %s: %v", field.QualifiedType, err)
	}

	if field.Repeated {
		typ = "[]" + typ
	}

	if field.Optional {
		typ = "*" + typ
	}

	return typ, nil
}

func getGoFieldName(name string) (string, string, error) {
	// Make sure the field name is snake case
	re := regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	if !re.MatchString(name) {
		return "", "", fmt.Errorf("%s is an invalid jrpc field name", name)
	}

	return snakeToCamelCase(name), name, nil
}

// getGoImportPath will return a go import path given a relative
// path to a def file and a package name that is relative to the
// def file location. E.g. if the module defined in the go.mod
// file is github.com/jakewright/home-automation, and a defPath
// of ../service.foo/foo.def and pkg of external are given,
// github.com/jakewright/home-automation/service.foo/external
// will be returned.
func getGoImportPath(defPath, pkg string) (string, error) {
	var module string
	var modFilePath string
	for i := 0; i < 10; i++ {
		modFilePath = strings.Repeat("../", i) + "go.mod"
		if i == 0 {
			modFilePath = "./" + modFilePath
		}

		b, err := ioutil.ReadFile(modFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}

		module = modfile.ModulePath(b)
		break
	}

	if module == "" {
		return "", fmt.Errorf("failed to find module path")
	}

	moduleRoot, err := filepath.Abs(filepath.Dir(modFilePath))
	if err != nil {
		return "", err
	}

	defPathAbs, err := filepath.Abs(defPath)
	if err != nil {
		return "", err
	}

	importPathRelToRoot, err := filepath.Rel(moduleRoot, defPathAbs)
	if err != nil {
		return "", err
	}

	svcImportPath := filepath.Dir(filepath.Join(module, importPathRelToRoot))

	return filepath.Join(svcImportPath, pkg), nil
}

// resolveTypeName will turn a fully-qualified type name into the go type and,
// if necessary, a path that needs to be imported.
func resolveTypeName(str, defPath string, imports map[string]*svcdef.Import, im *importManager) (string, error) {
	var goTypeName, importPath string
	var err error

	if strings.HasPrefix(str, ".") { // local type (message is defined in the same def file)
		goTypeName = strings.ReplaceAll(str[1:], ".", "_")

		// By convention, the type will be defined in the external package
		importPath, err = getGoImportPath(defPath, packageExternal)
		if err != nil {
			return "", err
		}

	} else if parts := strings.SplitN(str, ".", 2); len(parts) == 2 { // imported type
		goTypeName = strings.ReplaceAll(parts[1], ".", "_")

		// Expect to find an import with an alias of parts[0], and again, by
		// convention, the type name will be defined in the external package.
		importPath, err = getGoImportPath(imports[parts[0]].Path, packageExternal)
		if err != nil {
			return "", err
		}

	} else { // "built-in" type
		data, ok := typeMap[str]
		if !ok {
			return "", fmt.Errorf("invalid type %s", str)
		}

		goTypeName, importPath = data.GoType, data.ImportPath
	}

	alias := im.add(importPath)
	if alias != "" {
		goTypeName = alias + "." + goTypeName
	}

	return goTypeName, nil
}

type typeData struct {
	GoType     string
	ImportPath string
}

var typeMap = map[string]typeData{
	"bool":    {"bool", ""},
	"string":  {"string", ""},
	"int32":   {"int32", ""},
	"int64":   {"int64", ""},
	"uint32":  {"uint32", ""},
	"uint64":  {"uint64", ""},
	"float32": {"float32", ""},
	"float64": {"float64", ""},
	"bytes":   {"[]byte", ""},
	"time":    {"Time", "time"},
}

func snakeToCamelCase(s string) string {
	var camel string
	var upper bool

	for i, c := range s {
		switch {
		case c == '_':
			upper = true
		case i == 0, upper:
			camel += strings.ToUpper(string(c))
			upper = false
		default:
			camel += string(c)
		}
	}

	return camel
}
