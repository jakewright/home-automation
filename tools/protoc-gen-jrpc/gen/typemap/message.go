package typemap

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

// MessageDefinition represents a message defined in the proto
type MessageDefinition struct {
	// Descriptor is is the DescriptorProto defining the message.
	Descriptor *descriptor.DescriptorProto
	// File is the File that the message was defined in. Or, if it has been
	// publicly imported, what File was that import performed in?
	File *descriptor.FileDescriptorProto
	// Parent is the parent message, if this was defined as a nested message. If
	// this was defiend at the top level, parent is nil.
	Parent *MessageDefinition
	// Comments describes the comments surrounding a message's definition. If it
	// was publicly imported, then these comments are from the actual source file,
	// not the file that the import was performed in.
	Comments Comments

	// path is the 'SourceCodeInfo' path. See the documentation for
	// github.com/golang/protobuf/protoc-gen-go/descriptor.SourceCodeInfo for an
	// explanation of its format.
	path []int32
}

// ProtoName returns the dot-delimited, fully-qualified protobuf name of the message.
func (m *MessageDefinition) ProtoName() string {
	prefix := "."
	if pkg := m.File.GetPackage(); pkg != "" {
		prefix += pkg + "."
	}

	if lineage := m.Lineage(); len(lineage) > 0 {
		for _, parent := range lineage {
			prefix += parent.Descriptor.GetName() + "."
		}
	}

	return prefix + m.Descriptor.GetName()
}

// Lineage returns m's parental chain all the way back up to a top-level message
// definition. The first element of the returned slice is the highest-level parent.
func (m *MessageDefinition) Lineage() []*MessageDefinition {
	var parents []*MessageDefinition
	for p := m.Parent; p != nil; p = p.Parent {
		parents = append([]*MessageDefinition{p}, parents...)
	}
	return parents
}

// descendants returns all the submessages defined within m, and all the
// descendants of those, recursively.
func (m *MessageDefinition) descendants() []*MessageDefinition {
	descendants := make([]*MessageDefinition, 0)
	for i, child := range m.Descriptor.NestedType {
		path := append(m.path, []int32{messageMessagePath, int32(i)}...)
		childDef := &MessageDefinition{
			Descriptor: child,
			File:       m.File,
			Parent:     m,
			Comments:   commentsAtPath(path, m.File),
			path:       path,
		}
		descendants = append(descendants, childDef)
		descendants = append(descendants, childDef.descendants()...)
	}
	return descendants
}


