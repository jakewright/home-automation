package proto

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

type Proto struct {
	Files []*File
}

type File struct {
	Name string

	ImportPath string
	ProtoPackage string
	GoPackage    string

	Services []*Service
	Messages []*Message
}

type Service struct {
	Methods []*Method
}

type Method struct {
}

type Message struct {
	Name string

	// Path defines the location in the source file of this message
	Path []int32

	Parent *Message

	// ProtoName is the dot-delimited, fully-
	// qualified protobuf name of the message.
	ProtoName string
}

func Parse(fileDescriptors []*descriptor.FileDescriptorProto) (*Proto, error) {
	// Index all of the file descriptors by name
	fileDescriptorsByName := make(map[string]*descriptor.FileDescriptorProto)
	for _, d := range fileDescriptors {
		fileDescriptorsByName[d.GetName()] = d
	}

	for _, fd := range fileDescriptors {
		for _, md := range fd.MessageType {
			// Generate the ProtoName for this message
			prefix := "."
			if pkg := fd.GetPackage(); pkg != "" {
				prefix += pkg + "."
			}

		}

		//services := make([]*Service, len(d.Service))
		//for _, s := range d.Service {
		//
		//}

	}

	return &Proto{Files: files}, nil
}

func parseMessages(fileDescriptor *descriptor.FileDescriptorProto, fileDescriptorsByName map[string]*descriptor.FileDescriptorProto) {
	var messages []*Message

	for i, md := range fileDescriptor.MessageType {

		message := &Message{
			Path: []int32{messagePath, int32(i)},
		}
	}
}

func parseMessage(md *descriptor.DescriptorProto, parent *Message) []*Message {
	var path []int32

	message := &Message{
		Path:

	}
}

func parsePublicMessage() {
	panic("Parsing public messages is not implemented")
}
