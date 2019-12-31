package typemap

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Registry keeps maps of files by name and messages by name
type Registry struct {
	allFiles    []*descriptor.FileDescriptorProto
	filesByName map[string]*descriptor.FileDescriptorProto

	// Mapping of fully-qualified names to their definitions
	messagesByProtoName map[string]*MessageDefinition
}

// New returns a new typemap registry for the given set of files
func New(files []*descriptor.FileDescriptorProto) *Registry {
	r := &Registry{
		allFiles:            files,
		filesByName:         make(map[string]*descriptor.FileDescriptorProto),
		messagesByProtoName: make(map[string]*MessageDefinition),
	}

	// First, index the file descriptors by name. We need this so
	// messageDefsForFile can correctly scan imports.
	for _, f := range files {
		r.filesByName[f.GetName()] = f
	}

	// Next, index all the message definitions by their fully-qualified proto
	// names.
	for _, f := range files {
		defs := messageDefsForFile(f, r.filesByName)
		for name, def := range defs {
			r.messagesByProtoName[name] = def
		}
	}
	return r
}

// MethodInputDefinition returns the MessageDefinition of the input type of the given method
func (r *Registry) MethodInputDefinition(method *descriptor.MethodDescriptorProto) *MessageDefinition {
	return r.messagesByProtoName[method.GetInputType()]
}

// MethodOutputDefinition returns the MessageDefinition of the output type of the given method
func (r *Registry) MethodOutputDefinition(method *descriptor.MethodDescriptorProto) *MessageDefinition {
	return r.messagesByProtoName[method.GetOutputType()]
}

// MessageDefinition returns the MessageDefinition for the message with the given name
func (r *Registry) MessageDefinition(name string) *MessageDefinition {
	return r.messagesByProtoName[name]
}

// messageDefsForFile gathers a mapping of fully-qualified protobuf names to
// their definitions. It scans a singles file at a time. It requires a mapping
// of .proto file names to their definitions in order to correctly handle
// 'import public' declarations; this mapping should include all files
// transitively imported by f.
func messageDefsForFile(f *descriptor.FileDescriptorProto, filesByName map[string]*descriptor.FileDescriptorProto) map[string]*MessageDefinition {
	byProtoName := make(map[string]*MessageDefinition)
	// First, gather all the messages defined at the top level.
	for i, d := range f.MessageType {
		path := []int32{messagePath, int32(i)}
		def := &MessageDefinition{
			Descriptor: d,
			File:       f,
			Parent:     nil,
			Comments:   commentsAtPath(path, f),
			path:       path,
		}

		byProtoName[def.ProtoName()] = def
		// Next, all nested message definitions.
		for _, child := range def.descendants() {
			byProtoName[child.ProtoName()] = child
		}
	}

	// Finally, all messages imported publicly.
	for _, depIdx := range f.PublicDependency {
		depFileName := f.Dependency[depIdx]
		depFile := filesByName[depFileName]
		depDefs := messageDefsForFile(depFile, filesByName)
		for _, def := range depDefs {
			imported := &MessageDefinition{
				Descriptor: def.Descriptor,
				File:       f,
				Parent:     def.Parent,
				Comments:   commentsAtPath(def.path, depFile),
				path:       def.path,
			}
			byProtoName[imported.ProtoName()] = imported
		}
	}

	return byProtoName
}
