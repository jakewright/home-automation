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

	generators := []func(*protoparse.File) (string, []byte, error){
		generateRouter,
		generateValidate,
		generateFirehose,
	}

	for _, file := range files {
		for _, generate := range generators {
			name, b, err := generate(file)
			if err != nil {
				return nil, err
			}

			if len(b) == 0 {
				continue
			}

			// Format the code
			b, err = format.Source(b)
			if err != nil {
				return nil, err
			}

			filename := path.Dir(file.Name) + "/" + name + ".pb.go"

			responseFiles = append(responseFiles, &plugin_go.CodeGeneratorResponse_File{
				Name:    &filename,
				Content: proto.String(string(b)),
			})
		}
	}

	return &plugin_go.CodeGeneratorResponse{
		File: responseFiles,
	}, nil
}

func generateRouter(file *protoparse.File) (string, []byte, error) {
	if len(file.Services) != 1 {
		return "", nil, nil
	}

	data, err := createRouterTemplateData(file, file.Services[0])
	if err != nil {
		return "", nil, err
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := routerTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "router", buf.Bytes(), nil
}

func generateValidate(file *protoparse.File) (string, []byte, error) {
	data, err := createValidateTemplateData(file)
	if err != nil {
		return "", nil, err
	}

	if len(data.Messages) == 0 {
		return "", nil, nil
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := validateTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "validate", buf.Bytes(), nil
}

func generateFirehose(file *protoparse.File) (string, []byte, error) {
	data, err := createFirehoseTemplateData(file)
	if err != nil {
		return "", nil, err
	}

	if len(data.Events) == 0 {
		return "", nil, nil
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := firehoseTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "firehose", buf.Bytes(), nil
}
