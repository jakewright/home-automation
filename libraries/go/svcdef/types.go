package svcdef

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

	// Options are arbitrary options defined
	// at the top level in the def file
	Options map[string]interface{}
}

// AddMessage adds a new message to the file
func (f *File) AddMessage(msg *Message) {
	f.Messages = append(f.Messages, msg)
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

// Message is a representation of a message definition
type Message struct {
	// Name is the simple name given in the file, e.g. "Bar"
	Name string

	// QualifiedName is the fully-qualified name.
	// For a message defined at the top-level, this
	// will be the same as Name, e.g. "Bar". For
	// a nested message, it will be prefixed with
	// the parent's qualified name. E.g. "Foo.Bar".
	// Note that it does not have a leading dot.
	QualifiedName string

	// Fields are the type-name pairs defined in the message
	Fields []*Field

	// Nested is the list of nested messages defined
	// within this message
	Nested []*Message

	// Options are arbitrary options defined within the message
	Options map[string]interface{}
}

// AddMessage adds a new nested message
func (m *Message) AddMessage(msg *Message) {
	m.Nested = append(m.Nested, msg)
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
	// Name is the simple type name given in the def file
	Name string

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
