package main

import (
	"bytes"
	"fmt"
	"go/format"
	"path"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"

	"github.com/jakewright/home-automation/libraries/go/protoparse"
	jrpcproto "github.com/jakewright/home-automation/tools/protoc-gen-jrpc/proto"
)

func generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	files, err := protoparse.Parse(req)
	if err != nil {
		return nil, err
	}

	// Return silently if this isn't a valid service proto
	if len(files) != 1 {
		return nil, nil
	}
	file := files[0]
	if len(file.Services) != 1 {
		return nil, nil
	}
	service := file.Services[0]

	data, err := createTemplateData(file, service)
	if err != nil {
		return nil, err
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, data); err != nil {
		panic(err)
	}

	// Format the code
	b, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	// Construct the filename
	filename := file.Name
	if ext := path.Ext(filename); ext == ".proto" {
		filename = filename[:len(filename)-len(ext)]
	}
	filename += ".rpc.go"

	return &plugin_go.CodeGeneratorResponse{
		File: []*plugin_go.CodeGeneratorResponse_File{{
			Name:    &filename,
			Content: proto.String(string(b)),
		}},
	}, nil
}

func createTemplateData(file *protoparse.File, service *protoparse.Service) (*data, error) {
	// Get the service options
	opts, err := service.GetExtension(jrpcproto.E_Router)
	if err != nil {
		return nil, err
	}
	router := opts.(*jrpcproto.Router)

	methods := make([]*method, len(service.Methods))
	for i, m := range service.Methods {
		// Get the handler options
		opts, err := m.GetExtension(jrpcproto.E_Handler)
		if err != nil {
			return nil, err
		}
		handler := opts.(*jrpcproto.Handler)

		// Prepend the types with the package name if different from
		// the package name of the file we're generating
		inputType := m.InputType.GoTypeName
		if m.InputType.File.GoPackage != file.GoPackage {
			inputType = m.InputType.File.GoPackage + "." + inputType
		}
		outputType := m.OutputType.GoTypeName
		if m.OutputType.File.GoPackage != file.GoPackage {
			outputType = m.OutputType.File.GoPackage + "." + outputType
		}

		methods[i] = &method{
			Name:       m.Name,
			InputType:  inputType,
			OutputType: outputType,
			HTTPMethod: handler.Method,
			Path:       handler.Path,
			URL:        router.Name + handler.Path,
		}
	}

	imports := append(file.Imports,
		&protoparse.Import{Alias: "", Path: "github.com/jakewright/home-automation/libraries/go/request"},
		&protoparse.Import{Alias: "", Path: "github.com/jakewright/home-automation/libraries/go/response"},
		&protoparse.Import{Alias: "", Path: "github.com/jakewright/home-automation/libraries/go/router"},
		&protoparse.Import{Alias: "", Path: "github.com/jakewright/home-automation/libraries/go/rpc"},
		&protoparse.Import{Alias: "", Path: "github.com/jakewright/home-automation/libraries/go/slog"},
	)

	routerName, err := generateRouterName(service.Name)
	if err != nil {
		return nil, err
	}

	return &data{
		PackageName: file.GoPackage,
		RouterName:  routerName,
		Imports:     imports,
		Methods:     methods,
	}, nil
}

func generateRouterName(serviceName string) (string, error) {
	err := fmt.Errorf("service name should be alphanumeric camelcase ending with \"Service\"")

	r := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*Service$`)
	if ok := r.MatchString(serviceName); !ok {
		return "", err
	}

	return strings.Title(strings.TrimSuffix(serviceName, "Service")) + "Router", nil
}
