package svcdef

import "strings"

// File represents a .def file
type File struct {
	// Path is the path that was given to the Parse function.
	// It'll be relative to the current working directory from
	// where Parse() was called. This is true for imports as
	// well.
	Path string

	// Imports is a map of alias to import
	// e.g. in the statement import foo "../foo/foo.def",
	// foo is the alias
	Imports map[string]*Import

	// Service is a representation of the service
	// defined in the file. Only one service can
	// be defined in a single def file.
	Service *Service

	// Messages are all of the message types defined in the file
	Messages []*Message

	// FlatMessages contains all of the messages, including
	// nested and imported.
	FlatMessages []*Message

	// Options are arbitrary options defined
	// at the top level in the def file
	Options map[string]interface{}
}

func (f *File) addImport(alias string, imp *Import) {
	if f.Imports == nil {
		f.Imports = map[string]*Import{}
	}

	f.Imports[alias] = imp
}

func (f *File) addMessage(msg *Message) {
	f.Messages = append(f.Messages, msg)
}

func (f *File) addOption(key string, value interface{}) {
	if f.Options == nil {
		f.Options = map[string]interface{}{}
	}

	f.Options[key] = value
}

// Import represents an imported file
type Import struct {
	*File

	// Alias is the mandatory alias given to the import
	Alias string

	// Path is the relative path as written in the def file
	Path string
}

// Service is a representation of a service definition
type Service struct {
	// Name is the name given in the def file
	Name string

	// RPCs is the set of rpc statements in the service definition
	RPCs []*RPC

	// Options are arbitrary options defined within the service
	Options map[string]interface{}
}

func (s *Service) addOption(key string, value interface{}) {
	if s.Options == nil {
		s.Options = map[string]interface{}{}
	}

	s.Options[key] = value
}

// RPC is a representation of an rpc definition
type RPC struct {
	// Name is the name given in the def file
	Name string

	// InputType is the input type
	InputType *Type

	// OutputType is the output type
	OutputType *Type

	// Options are arbitrary options defined within the RPC
	Options map[string]interface{}
}

func (r *RPC) addOption(key string, value interface{}) {
	if r.Options == nil {
		r.Options = map[string]interface{}{}
	}

	r.Options[key] = value
}

// Message is a representation of a message definition
type Message struct {
	// Name is the simple name given in the file, e.g. "Bar"
	Name string

	// QualifiedName is the fully-qualified name.
	// For a message defined in the main def file,
	// it will be the name prefixed with a dot, e.g.
	// ".Bar". For a nested message, it will also
	// include the parent's lineage, e.g. ".Foo.Bar".
	// For an imported message, it will be prefixed
	// with the import alias, e.g. "Foo.Bar".
	QualifiedName string

	// Fields are the type-name pairs defined in the message
	Fields []*Field

	// Nested is the list of nested messages defined
	// within this message
	Nested []*Message

	// Options are arbitrary options defined within the message
	Options map[string]interface{}
}

func (m *Message) addField(f *Field) {
	m.Fields = append(m.Fields, f)
}

func (m *Message) addMessage(msg *Message) {
	m.Nested = append(m.Nested, msg)
}

func (m *Message) addOption(key string, value interface{}) {
	if m.Options == nil {
		m.Options = map[string]interface{}{}
	}

	m.Options[key] = value
}

// Lineage returns the file alias and a slice of name parts.
// If the message was defined in the main def file, the first
// return value will be the empty string.
func (m *Message) Lineage() (string, []string) {
	parts := strings.Split(m.QualifiedName, ".")
	return parts[0], parts[1:]
}

// Field is a representation of a type-name pair
type Field struct {
	// name is the name given in the def file
	Name string

	// Type is the type of the field
	Type *Type

	// Options are arbitrary options defined in parenthesis after the field definition
	Options map[string]interface{}
}

// Type is a representation of a type
type Type struct {
	// Name is the simple name of the type e.g. "map",
	// Note that in the case of repeated types e.g. "[]int"
	// the Name is just "int" and Repeated is set to true.
	Name string

	// Original is the original type string from the def file
	// e.g. "map[string]string", "[]string"
	Original string

	// Qualified is the fully-qualified type name.
	// It can be
	//   - the name of a message defined in this file
	//     - e.g. ".Bar", or ".Bar.Baz"
	//   - the name of a message defined in another file
	//     - e.g. "foo.Bar"
	//   - a custom type for the code generator to parse
	//     - e.g. "string"
	//     - in the case of a map e.g. map[string]int
	//       the generator should use the Map* fields to
	//       understand how to interpret the type
	Qualified string

	// Message is the message that the type refers to.
	// It is nil if it is a "simple", non-message type.
	//Message *Message

	// Repeated is set if [] appears before the type name
	Repeated bool

	// Optional is set if * appears before the type name
	Optional bool

	// Map is true if the type is a map
	Map bool

	// MapKey is the type of the map's keys
	MapKey *Type

	// MapValue is the type of the map's values
	MapValue *Type
}
