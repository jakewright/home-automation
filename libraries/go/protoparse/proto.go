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

	PackageComments *Comments
	Imports         []*Import

	Services []*Service
	Messages []*Message

	descriptor *descriptor.FileDescriptorProto
}

func (f *File) GetPackageComments() *Comments {
	if f.PackageComments != nil {
		return f.PackageComments
	}

	return &Comments{}
}

// Service represents a service definition in a proto file
type Service struct {
	Name    string
	Methods []*Method

	// Comments defines the comments attached to the message
	Comments *Comments

	path       []int32
	descriptor *descriptor.ServiceDescriptorProto
}

// GetExtension can be used to get custom options set on a service
func (s *Service) GetExtension(extension *proto.ExtensionDesc) (interface{}, error) {
	return proto.GetExtension(s.descriptor, extension)
}

func (s *Service) GetComments() *Comments {
	if s.Comments != nil {
		return s.Comments
	}

	return &Comments{}
}

// Method represents a method defined in a service
type Method struct {
	Name       string
	InputType  *Message
	OutputType *Message

	// Comments defines the comments attached to the method
	Comments *Comments

	path       []int32
	descriptor *descriptor.MethodDescriptorProto
}

// GetExtension can be used to get custom options set on a method
func (m *Method) GetExtension(extension *proto.ExtensionDesc) (interface{}, error) {
	return proto.GetExtension(m.descriptor, extension)
}

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

	// File is the file in which the message was defined (or the file
	// in which it was imported in the case of a public import).
	File *File

	// Parent is set if the message definition was nested inside another
	Parent *Message

	path       []int32
	descriptor *descriptor.DescriptorProto
}

func (m *Message) GetComments() *Comments {
	if m.Comments != nil {
		return m.Comments
	}

	return &Comments{}
}

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
	Alias string
	Path  string
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
