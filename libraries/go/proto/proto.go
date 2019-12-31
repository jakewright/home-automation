package proto

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
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

type Proto struct {
	Files []*File
}

type File struct {
	// Name is the TODO
	Name string

	ImportPath      string
	ProtoPackage    string
	GoPackage       string
	PackageComments *Comments
	Imports         []*Import

	Services []*Service
	Messages []*Message

	descriptor *descriptor.FileDescriptorProto
}

type Service struct {
	Name    string
	Methods []*Method

	// Comments defines the comments attached to the message
	Comments *Comments

	path       []int32
	descriptor *descriptor.ServiceDescriptorProto
}

func (s *Service) GetExtension(extension *proto.ExtensionDesc) (interface{}, error) {
	return proto.GetExtension(s.descriptor, extension)
}

type Method struct {
	Name       string
	InputType  *Message
	OutputType *Message

	// Comments defines the comments attached to the method
	Comments *Comments

	path       []int32
	descriptor *descriptor.MethodDescriptorProto
}

func (m *Method) GetExtension(extension *proto.ExtensionDesc) (interface{}, error) {
	return proto.GetExtension(m.descriptor, extension)
}

type Message struct {
	// Name is the simple name of the message
	Name string

	// ProtoName is the dot-delimited, fully-
	// qualified protobuf name of the message.
	ProtoName string

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

type Comments struct {
	Leading         string
	LeadingDetached []string
	Trailing        string
}

type Import struct {
	Alias string
	Path  string
}

func Parse(fileDescriptors []*descriptor.FileDescriptorProto) (*Proto, error) {
	// Create a slice of Files. For now, services and messages
	// will be empty but we will populate them later.
	files := make([]*File, len(fileDescriptors))
	for i, fd := range fileDescriptors {
		importPath, goPackage, ok := readGoPackageOption(fd)
		if !ok {
			return nil, fmt.Errorf("go_package option is not set in file %s", fd.GetName())
		}

		comments := commentsAtPath([]int32{packagePath}, fd)

		files[i] = &File{
			Name:            fd.GetName(),
			ProtoPackage:    fd.GetPackage(),
			ImportPath:      importPath,
			GoPackage:       goPackage,
			PackageComments: comments,
			descriptor:      fd,
		}
	}

	// Index all of the files by name
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

	return &Proto{Files: files}, nil
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

	// Create the fully-qualified proto name
	protoName := "." + md.GetName()
	for p := parent; p != nil; p = p.Parent {
		protoName = "." + p.Name + protoName
	}
	if pkg := f.ProtoPackage; pkg != "" {
		protoName = "." + pkg + protoName
	}

	message := &Message{
		Name:       md.GetName(),
		Parent:     parent,
		ProtoName:  protoName,
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
	path := importFile.ImportPath

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
		pathEqual(path, l.Path)
		{
			return &Comments{
				Leading:         l.GetLeadingComments(),
				LeadingDetached: l.GetLeadingDetachedComments(),
				Trailing:        l.GetTrailingComments(),
			}
		}
	}

	return &Comments{}
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
