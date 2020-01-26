package main

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/tools/imports"

	"github.com/jakewright/home-automation/libraries/go/svcdef"
)

const (
	packageExternal = "external"
	packageRouter   = "handler"
)

type options struct {
	DefPath           string
	RouterPackageName string
}

func generate(defPath string, file *svcdef.File) error {
	generators := []func(*options, *svcdef.File) (string, []byte, error){
		generateTypes,
		generateRouter,
		generateFirehose,
	}

	opts := &options{
		DefPath: defPath,
	}

	for _, generate := range generators {
		filename, b, err := generate(opts, file)
		if err != nil {
			panic(err)
		}

		b, err = imports.Process(filename, b, &imports.Options{
			Comments: true,
		})
		if err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(filename, b, 0644); err != nil {
			panic(err)
		}
	}

	return nil
}

func generateTypes(opts *options, file *svcdef.File) (string, []byte, error) {
	data, err := createTypesTemplateData(opts, file)
	if err != nil {
		return "", nil, err
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := typesTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "./" + data.PackageName + "/types.go", buf.Bytes(), nil
}

func generateRouter(opts *options, file *svcdef.File) (string, []byte, error) {
	data, err := createRouterTemplateData(opts, file)
	if err != nil {
		return "", nil, err
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := routerTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "./" + data.PackageName + "/gen.go", buf.Bytes(), nil
}

func generateFirehose(opts *options, file *svcdef.File) (string, []byte, error) {
	data, err := createFirehoseTemplateData(opts, file)
	if err != nil {
		return "", nil, err
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := firehoseTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "./" + data.PackageName + "/firehose.go", buf.Bytes(), nil
}
