package main

import (
	"bytes"
	"go/format"
	"path"

	"github.com/golang/protobuf/proto"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"

	"github.com/jakewright/home-automation/libraries/go/protoparse"
)

func generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	files, err := protoparse.Parse(req)
	if err != nil {
		return nil, err
	}

	var responseFiles []*plugin_go.CodeGeneratorResponse_File

	for _, file := range files {
		if len(file.Services) == 1 {
			responseFile, err := generateRouterFile(file, file.Services[0])
			if err != nil {
				return nil, err
			}

			responseFiles = append(responseFiles, responseFile)
		}

		responseFile, err := generateValidateFile(file)
		if err != nil {
			return nil, err
		}

		responseFiles = append(responseFiles, responseFile)
	}

	return &plugin_go.CodeGeneratorResponse{
		File: responseFiles,
	}, nil
}

func generateRouterFile(file *protoparse.File, service *protoparse.Service) (*plugin_go.CodeGeneratorResponse_File, error) {
	data, err := createRouterTemplateData(file, service)
	if err != nil {
		return nil, err
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := routerTemplate.Execute(buf, data); err != nil {
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

	return &plugin_go.CodeGeneratorResponse_File{
		Name:    &filename,
		Content: proto.String(string(b)),
	}, nil
}

func generateValidateFile(file *protoparse.File) (*plugin_go.CodeGeneratorResponse_File, error) {
	data, err := createValidateTemplateData(file)
	if err != nil {
		return nil, err
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := validateTemplate.Execute(buf, data); err != nil {
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
	filename += ".validate.go"

	return &plugin_go.CodeGeneratorResponse_File{
		Name:    &filename,
		Content: proto.String(string(b)),
	}, nil
}
