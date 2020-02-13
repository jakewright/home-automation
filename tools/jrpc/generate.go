package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"

	"github.com/jakewright/home-automation/libraries/go/svcdef"
)

const (
	packageDirExternal = "def"
	packageDirRouter   = "handler"
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

		if filename == "" || b == nil {
			continue
		}

		b, err = imports.Process(filename, b, &imports.Options{
			Comments: true,
		})
		if err != nil {
			panic(err)
		}

		dir := filepath.Dir(filename)

		// Create the directory if necessary
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0700); err != nil {
				panic(err)
			}
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

	if data == nil {
		return "", nil, nil
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := typesTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "./" + data.PackageDir + "/types.go", buf.Bytes(), nil
}

func generateRouter(opts *options, file *svcdef.File) (string, []byte, error) {
	data, err := createRouterTemplateData(opts, file)
	if err != nil {
		return "", nil, err
	}

	if data == nil {
		return "", nil, nil
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := routerTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "./" + data.PackageDir + "/gen.go", buf.Bytes(), nil
}

func generateFirehose(opts *options, file *svcdef.File) (string, []byte, error) {
	data, err := createFirehoseTemplateData(opts, file)
	if err != nil {
		return "", nil, err
	}

	if data == nil {
		return "", nil, nil
	}

	// Generate the code
	buf := &bytes.Buffer{}
	if err := firehoseTemplate.Execute(buf, data); err != nil {
		return "", nil, err
	}

	return "./" + data.PackageDir + "/firehose.go", buf.Bytes(), nil
}
