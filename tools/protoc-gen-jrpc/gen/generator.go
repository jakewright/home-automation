package gen

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/jakewright/home-automation/tools/protoc-gen-jrpc/gen/typemap"
	"github.com/jakewright/home-automation/tools/protoc-gen-jrpc/jrpc"
	jrpcproto "github.com/jakewright/home-automation/tools/protoc-gen-jrpc/proto"
)

func Generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	if len(req.FileToGenerate) != 1 {
		err := "JRPC only supports single proto files per service"
		return &plugin.CodeGeneratorResponse{Error: &err}
	}

	// Iterate over all proto descriptors to find the one
	// that we're actually supposed to generate code for.
	// The others are just the things it imports.
	var protoFileToGenerate *descriptor.FileDescriptorProto
	allProtoFiles := req.ProtoFile
	for _, file := range allProtoFiles {
		if file.GetName() == req.FileToGenerate[0] {
			protoFileToGenerate = file
			break
		}
	}
	if protoFileToGenerate == nil {
		panic("protoFileToGenerate should not be nil at this point")
	}
}

func templateDataFromFile(file *descriptor.FileDescriptorProto, files []*descriptor.FileDescriptorProto) (*templateData, error) {
	if len(file.Service) < 1 {
		stderr("No services in file %s; skipping...", file.GetName())
		return nil, nil
	}
	if len(file.Service) > 1 {
		stderr("Too many services in file %s; skipping...", file.GetName())
		return nil, nil
	}

	registry := typemap.New(files)

	registry.MethodInputDefinition()

	packageName, err := goPackageName(file)
	if err != nil {
		return nil, err
	}

	packageComment, err :=

	return &templateData {
		PackageName: packageName,
	}, nil
}

type Generator struct {
	reg *typemap.Registry
	pkg string
}

func (g *Generator) Generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	if len(req.FileToGenerate) != 1 {
		panic("JRPC only supports single proto files per service")
	}

	// Iterate over all proto descriptors to find the one
	// that we're actually supposed to generate code for.
	// The others are just the things it imports.
	var fileToGenerate *descriptor.FileDescriptorProto
	for _, file := range req.ProtoFile {
		if file.GetName() == req.FileToGenerate[0] {
			fileToGenerate = file
			break
		}
	}

	rsp := &plugin.CodeGeneratorResponse{}

	jFile := g.generate(fileToGenerate, req.ProtoFile)
	if jFile != nil {
		rsp.File = append(rsp.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(goFileName(fileToGenerate)),
			Content: proto.String(jFile.Generate()),
		})
	}

	return rsp
}

// generate converts a proto file descriptor into a JRPC file descriptor
func (g *Generator) generate(file *descriptor.FileDescriptorProto, files []*descriptor.FileDescriptorProto) *jrpc.FileDescriptor {
	if len(file.Service) < 1 {
		stderr("No services in file %s; skipping...", file.GetName())
		return nil
	}
	if len(file.Service) > 1 {
		stderr("Too many services in file %s; skipping...", file.GetName())
		return nil
	}

	g.reg = typemap.New(files)
	g.pkg = goPackageName(file)

	jFile := jrpc.New()

	jFile.PackageName = g.pkg
	jFile.PackageComment = g.generatePackageComment(file)
	jFile.Imports = append(jFile.Imports, g.generateImports(file)...)
	jFile.Service = &jrpc.Service{}

	opts, err := proto.GetExtension(file.Service[0].Options, jrpcproto.E_Router)
	if err != nil {
		panic(err)
	}
	router := opts.(*jrpcproto.Router)

	for _, m := range file.Service[0].Method {
		opts, err := proto.GetExtension(m.Options, jrpcproto.E_Handler)
		if err != nil {
			panic(err)
		}
		handler := opts.(*jrpcproto.Handler)

		rpc := &jrpc.RPC{
			Name:         m.GetName(),
			URL:          router.Name + handler.Path,
			HTTPMethod:   handler.Method,
			RequestType:  g.goTypeName(g.reg.MethodInputDefinition(m)),
			ResponseType: g.goTypeName(g.reg.MethodOutputDefinition(m)),
		}

		jFile.Service.RPCs = append(jFile.Service.RPCs, rpc)
	}

	return jFile
}

func goPackageComment(registry *typemap.Registry, file *descriptor.FileDescriptorProto) ([]string, error) {
	fileComments, err := registry.FileComments(file)
	if err != nil {
		return nil, err
	}

	var packageComment []string

	if fileComments.Leading != "" {
		for _, line := range strings.Split(fileComments.Leading, "\n") {
			line = strings.TrimPrefix(line, " ")
			if line == "" {
				continue
			}
			packageComment = append(packageComment, line)
		}
	}

	return packageComment, nil
}

func (g *Generator) generateImports(file *descriptor.FileDescriptorProto) []*jrpc.Import {
	var imports []*jrpc.Import

	// It's legal to import output definitions from other proto files
	// JRPC doesn't support importing input definitions because it needs
	// to add a function with the input type as the receiver.
	for _, service := range file.Service {
		for _, method := range service.Method {
			if g.reg.MethodInputDefinition(method).File != file {
				panic(fmt.Sprintf("JPRC does not support imported input definitions: %s", method.GetName()))
			}

			def := g.reg.MethodOutputDefinition(method)
			if def.File != file { // Don't need to import the current file
				path, pkg, ok := goPackageOption(def.File)
				if !ok {
					panic(fmt.Sprintf("File %s does not have option go_package set", file.GetName()))
				}

				imports = append(imports, &jrpc.Import{
					Alias: pkg,
					Path:  path,
				})
			}
		}
	}

	return imports
}
