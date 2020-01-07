package protoparse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
)

const (
	// Tag numbers in FileDescriptorProto
	packagePath = 2 // package
	messagePath = 4 // message_type
	enumPath    = 5 // enum_type
	servicePath = 6 // service

	// Tag numbers in DescriptorProto
	messageFieldPath   = 2 // field
	messageMessagePath = 3 // nested_type
	messageEnumPath    = 4 // enum_type
	messageOneofPath   = 8 // oneof_decl

	// Tag numbers in ServiceDescriptorProto
	serviceNamePath    = 1 // name
	serviceMethodPath  = 2 // method
	serviceOptionsPath = 3 // options

	// Tag numbers in MethodDescriptorProto
	methodNamePath   = 1 // name
	methodInputPath  = 2 // input_type
	methodOutputPath = 3 // output_type
)

// File represents a parsed proto file
type File struct {
	// Name is the relative path and filename of the proto file
	Name string

	// ProtoPackage is the name defined in the proto's package statement
	ProtoPackage string

	// GoImportPath is the go import path of the package as defined by the go_package option
	GoImportPath string

	// GoPackage is the name defined by the go_package option
	GoPackage string

	// PackageComments are comments on the package definition
	PackageComments *Comments

	// Imports is the list of go imports this file needs for its service types
	Imports []*Import

	// Services are the services defined in the proto file
	Services []*Service

	// Messages are the messages defined in the proto
	// file and any messages from public imports
	Messages []*Message

	descriptor *descriptor.FileDescriptorProto
}

// GetPackageComments is a nil-safe getter for PackageComments
func (f *File) GetPackageComments() *Comments {
	if f.PackageComments != nil {
		return f.PackageComments
	}

	return &Comments{}
}

// Service represents a service definition in a proto file
type Service struct {
	// Name is the simple name of the service
	Name string

	// Methods are the methods (RPCs) defined in this service
	Methods []*Method

	// Comments defines the comments attached to the service
	Comments *Comments

	path       []int32
	descriptor *descriptor.ServiceDescriptorProto
}

// GetExtension can be used to get custom options set on a service
func (s *Service) GetExtension(extension *proto.ExtensionDesc) (interface{}, error) {
	return proto.GetExtension(s.descriptor.Options, extension)
}

// GetComments is a nil-safe getter for Comments
func (s *Service) GetComments() *Comments {
	if s.Comments != nil {
		return s.Comments
	}

	return &Comments{}
}

// Method represents a method defined in a service
type Method struct {
	// Name is the simple name for this method
	Name string

	// InputType is the message defined as the input type
	InputType *Message

	// OutputType is the message defined as the output type
	OutputType *Message

	// Comments defines the comments attached to the method
	Comments *Comments

	path       []int32
	descriptor *descriptor.MethodDescriptorProto
}

// GetExtension can be used to get custom options set on a method
func (m *Method) GetExtension(extension *proto.ExtensionDesc) (interface{}, error) {
	return proto.GetExtension(m.descriptor.Options, extension)
}

// GetComments is a nil-safe getter for Comments
func (m *Method) GetComments() *Comments {
	if m.Comments != nil {
		return m.Comments
	}

	return &Comments{}
}

// Message represents a message type defined in a proto file
type Message struct {
	// Name is the simple name of the message
	Name string

	// ProtoName is the dot-delimited, fully-
	// qualified protobuf name of the message.
	ProtoName string

	// GoTypeName is the name of the go type generated for this message
	GoTypeName string

	// Comments defines the comments attached to the message
	Comments *Comments

	// Fields are the fields defined in the message
	Fields []*Field

	// File is the file in which the message was defined (or the file
	// in which it was imported in the case of a public import).
	File *File

	// Parent is set if the message definition was nested inside another
	Parent *Message

	path       []int32
	descriptor *descriptor.DescriptorProto
}

// GetExtension can be used to get custom options set on a message
func (m *Message) GetExtension(extension *proto.ExtensionDesc) (interface{}, error) {
	if m.descriptor.Options == nil {
		return extension.ExtensionType, nil
	}

	return proto.GetExtension(m.descriptor.Options, extension)
}

// Field represents a field within a message
type Field struct {
	// Name is the name from the proto file e.g. email_address
	Name string

	// Repeated is true if this field has the repeated label
	Repeated bool

	// GoName is the name the field will get in the go struct e.g. EmailAddress
	GoName string

	// TypeName is the name of the type e.g. TYPE_STRING
	TypeName string

	// Type is the message if TypeName == TYPE_MESSAGE
	Type *Message

	descriptor *descriptor.FieldDescriptorProto
}

// GetExtension can be used to get custom options set on a field
func (f *Field) GetExtension(extension *proto.ExtensionDesc) (interface{}, error) {
	if f.descriptor.Options == nil {
		return extension.ExtensionType, nil
	}

	return proto.GetExtension(f.descriptor.Options, extension)
}

// GetComments is a nil-safe getter for Comments
func (m *Message) GetComments() *Comments {
	if m.Comments != nil {
		return m.Comments
	}

	return &Comments{}
}

// GetParent is a nil-safe getter for Parent
func (m *Message) GetParent() *Message {
	if m.Parent != nil {
		return m.Parent
	}

	return &Message{}
}

// Comments holds the set of comments associated with an entity
type Comments struct {
	// Leading are the comment lines directly above the line of code
	Leading []string

	// Leading detached are comment lines above
	// the line of code but not directly touching
	LeadingDetached [][]string

	// Trailing are comment lines directly under the line of code
	Trailing []string
}

// Import describes a go import
type Import struct {
	// Alias is the import alias if one is required
	Alias string

	// Path is the path to be imported
	Path string
}

// Parse turns a CodeGeneratorRequest into a parsed set of Files
func Parse(req *plugin_go.CodeGeneratorRequest) ([]*File, error) {
	// Create a slice of Files. For now, services and messages
	// will be empty but we will populate them later.
	files := make([]*File, len(req.ProtoFile))
	for i, fd := range req.ProtoFile {
		goImportPath, goPackage, ok := readGoPackageOption(fd)
		if !ok {
			return nil, fmt.Errorf("go_package option is not set in file %s", fd.GetName())
		}

		comments := commentsAtPath([]int32{packagePath}, fd)

		files[i] = &File{
			Name:            fd.GetName(),
			ProtoPackage:    fd.GetPackage(),
			GoImportPath:    goImportPath,
			GoPackage:       goPackage,
			PackageComments: comments,
			descriptor:      fd,
		}
	}

	// Index all of the files by name. This is needed so we
	// can correctly handle public imports when parsing messages.
	filesByName := make(map[string]*File)
	for _, f := range files {
		filesByName[f.Name] = f
	}

	// Parse the messages defined in each file
	messagesByProtoName := map[string]*Message{}
	for _, f := range files {
		f.Messages = parseMessages(f, filesByName)
		for _, msg := range f.Messages {
			messagesByProtoName[msg.ProtoName] = msg
		}
	}

	if err := parseMessageFields(messagesByProtoName); err != nil {
		return nil, err
	}

	// Parse services defined in each file
	for _, f := range files {
		f.Services, f.Imports = parseServices(f, messagesByProtoName)
	}

	// Pull out the files to generate. The rest are just imports.
	var filesToGenerate []*File
	for _, f := range files {
		for _, name := range req.FileToGenerate {
			if f.Name == name {
				filesToGenerate = append(filesToGenerate, f)
			}
		}
	}

	return filesToGenerate, nil
}

func readGoPackageOption(fd *descriptor.FileDescriptorProto) (string, string, bool) {
	opt := fd.GetOptions().GetGoPackage()
	if opt == "" {
		return "", "", false
	}

	// A semicolon-delimited suffix delimits the import path and package name.
	sc := strings.Index(opt, ";")
	if sc >= 0 {
		return opt[:sc], cleanPackageName(opt[sc+1:]), true
	}

	// The presence of a slash implies there's an import path
	slash := strings.LastIndex(opt, "/")
	if slash < 0 {
		return "", opt, true
	}

	return opt, opt[slash+1:], true
}

func cleanPackageName(name string) string {
	name = strings.Map(badToUnderscore, name)
	// Identifier must not be keyword
	if isGoKeyword(name) {
		name = "_" + name
	}
	// Identifier must not begin with digit
	if r, _ := utf8.DecodeRuneInString(name); unicode.IsDigit(r) {
		name = "_" + name
	}
	return name
}

func badToUnderscore(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return r
	}
	return '_'
}

func parseMessages(f *File, filesByName map[string]*File) []*Message {
	var messages []*Message

	// For each message defined in the proto file
	for i, md := range f.descriptor.MessageType {
		messages = append(messages, parseMessage(f, md, i, nil)...)
	}

	// For each public dependency where i is the dependency's index
	for _, i := range f.descriptor.PublicDependency {
		filename := f.descriptor.Dependency[i]
		depFd := filesByName[filename]
		importedMessages := parseMessages(depFd, filesByName)

		// For each publicly imported message, update the file to be f
		// instead of the file that the message was actually defined in.
		for _, md := range importedMessages {
			md.File = f
		}

		messages = append(messages, importedMessages...)
	}

	return messages
}

func parseMessage(f *File, md *descriptor.DescriptorProto, index int, parent *Message) []*Message {
	var path []int32
	if parent == nil {
		path = []int32{messagePath, int32(index)}
	} else {
		path = append(parent.path, messageMessagePath, int32(index))
	}

	// Create the fully-qualified proto name and the go type name.
	// TODO if the message names are not camelcase in the proto file,
	// the go type name will not necessarily come out right.
	protoName := "." + md.GetName()
	goTypeName := md.GetName()
	for p := parent; p != nil; p = p.Parent {
		protoName = "." + p.Name + protoName
		goTypeName = p.Name + "_" + goTypeName
	}
	if pkg := f.ProtoPackage; pkg != "" {
		protoName = "." + pkg + protoName
	}

	message := &Message{
		Name:       md.GetName(),
		File:       f,
		Parent:     parent,
		ProtoName:  protoName,
		GoTypeName: goTypeName,
		Comments:   commentsAtPath(path, f.descriptor),
		path:       path,
		descriptor: md,
	}

	// Get all nested definitions
	messages := []*Message{message}
	for i, child := range md.NestedType {
		messages = append(messages, parseMessage(f, child, i, message)...)
	}

	return messages
}

// parseMessageFields has to be done after parsing the messages themselves, because
// fields can have types that are other messages, so we need to know about all
// messages that exist before we can parse the fields.
func parseMessageFields(messagesByProtoName map[string]*Message) error {
	for _, m := range messagesByProtoName {
		fields := make([]*Field, len(m.descriptor.Field))
		for i, fd := range m.descriptor.Field {
			// The camelCase function is taken from the protobuf
			// package so this should be roughly accurate
			goName := camelCase(fd.GetName())

			// Get the type name. Note that for message types
			// this will be TYPE_MESSAGE. For the message's
			// actual name, check field.Type (set below).
			typeIndex := int32(fd.GetType())
			typeName := descriptor.FieldDescriptorProto_Type_name[typeIndex]

			var message *Message
			// If this is a message type, find the message.
			if fd.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
				var ok bool
				if message, ok = messagesByProtoName[fd.GetTypeName()]; !ok {
					return fmt.Errorf("field %s has type %s but could not find a message with that name", fd.GetName(), fd.GetTypeName())
				}
			}

			fields[i] = &Field{
				Name:       fd.GetName(),
				Repeated:   fd.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
				GoName:     goName,
				TypeName:   typeName,
				Type:       message,
				descriptor: fd,
			}
		}

		m.Fields = fields
	}

	return nil
}

func parseServices(f *File, messagesByProtoName map[string]*Message) ([]*Service, []*Import) {
	services := make([]*Service, len(f.descriptor.Service))
	var imports []*Import

	for i, sd := range f.descriptor.Service {
		path := []int32{servicePath, int32(i)}

		svc := &Service{
			Name:       sd.GetName(),
			Comments:   commentsAtPath(path, f.descriptor),
			path:       path,
			descriptor: sd,
		}

		methods, imps := parseMethods(f, svc, messagesByProtoName)
		svc.Methods = methods
		services[i] = svc
		imports = append(imports, imps...)
	}

	return services, imports
}

func parseMethods(f *File, s *Service, messagesByProtoName map[string]*Message) ([]*Method, []*Import) {
	methods := make([]*Method, len(s.descriptor.Method))
	var imports []*Import

	for i, md := range s.descriptor.Method {
		path := append(s.path, serviceMethodPath, int32(i))
		inputType := messagesByProtoName[md.GetInputType()]
		outputType := messagesByProtoName[md.GetOutputType()]

		methods[i] = &Method{
			Name:       md.GetName(),
			Comments:   commentsAtPath(path, f.descriptor),
			InputType:  inputType,
			OutputType: outputType,
			path:       path,
			descriptor: md,
		}

		if imp := createImport(inputType.File, f); imp != nil {
			imports = append(imports, imp)
		}
		if imp := createImport(outputType.File, f); imp != nil {
			imports = append(imports, imp)
		}
	}

	return methods, imports
}

func createImport(importFile, dstFile *File) *Import {
	// Files don't need to import themselves
	if importFile == dstFile {
		return nil
	}

	var alias string
	path := importFile.GoImportPath

	// If the go package is different to the import path, we need an alias.
	if !strings.HasSuffix(path, "/"+importFile.GoPackage) {
		alias = importFile.GoPackage
	}

	return &Import{
		Alias: alias,
		Path:  path,
	}
}

func commentsAtPath(path []int32, fd *descriptor.FileDescriptorProto) *Comments {
	if fd.SourceCodeInfo == nil {
		return &Comments{}
	}

	for _, l := range fd.SourceCodeInfo.Location {
		if pathEqual(path, l.Path) {
			var leadingDetached [][]string
			for _, s := range l.GetLeadingDetachedComments() {
				leadingDetached = append(leadingDetached, normalizeComment(s))
			}

			return &Comments{
				Leading:         normalizeComment(l.GetLeadingComments()),
				LeadingDetached: leadingDetached,
				Trailing:        normalizeComment(l.GetTrailingComments()),
			}
		}
	}

	return &Comments{}
}

// normalizeComment splits a comment line
// by \n characters and trims the leading space
func normalizeComment(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimPrefix(line, " ")
		lines = append(lines, line)
	}
	return lines
}

func pathEqual(path1, path2 []int32) bool {
	if len(path1) != len(path2) {
		return false
	}
	for i, v := range path1 {
		if path2[i] != v {
			return false
		}
	}
	return true
}

func isGoKeyword(s string) bool {
	return map[string]bool{
		"break":       true,
		"case":        true,
		"chan":        true,
		"const":       true,
		"continue":    true,
		"default":     true,
		"else":        true,
		"defer":       true,
		"fallthrough": true,
		"for":         true,
		"func":        true,
		"go":          true,
		"goto":        true,
		"if":          true,
		"import":      true,
		"interface":   true,
		"map":         true,
		"package":     true,
		"range":       true,
		"return":      true,
		"select":      true,
		"struct":      true,
		"switch":      true,
		"type":        true,
		"var":         true,
	}[s]
}

// camelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
