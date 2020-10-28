package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/imports"

	"github.com/jakewright/home-automation/libraries/go/svcdef"
	jrpcimports "github.com/jakewright/home-automation/tools/libraries/imports"
)

var resolver *jrpcimports.Resolver

func init() {
	var err error
	resolver, err = jrpcimports.NewResolver()
	if err != nil {
		panic(err)
	}
}

type generator interface {
	Init(*options, *svcdef.File)
	PackageDir() string
	Data(*jrpcimports.Manager) (interface{}, error)
	Template() (*template.Template, error)
	Filename() string
}

type baseGenerator struct {
	options *options
	file    *svcdef.File
}

func (g *baseGenerator) Init(options *options, file *svcdef.File) {
	g.options = options
	g.file = file
}

type options struct {
	DefPath           string
	RouterPackageName string
}

func generate(defPath string, file *svcdef.File) error {
	generators := []generator{
		&clientGenerator{},
		&firehoseGenerator{},
		&routerGenerator{},
		&typesGenerator{},
	}

	opts := &options{
		DefPath: defPath,
	}

	for _, generator := range generators {
		// Initialise the generator
		generator.Init(opts, file)

		// Generate the package directory name e.g. "routes"
		packageDir := generator.PackageDir()

		// Generate the filename
		filename := filepath.Join(filepath.Dir(defPath), packageDir, generator.Filename())

		// Get the full go import path of the package we're generating
		self, err := resolver.Resolve(file.Path, packageDir)
		if err != nil {
			return err
		}

		// Create an import manager
		im := jrpcimports.NewManager(self)

		// Generate the template data
		data, err := generator.Data(im)
		if err != nil {
			return err
		}
		// If there's nothing to generate, delete the file
		// in case it previously existed.
		if data == nil {
			if err := os.Remove(filename); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return err
			}

			continue
		}

		// Get the template
		tmpl, err := generator.Template()
		if err != nil {
			return err
		}

		// Generate the code
		buf := &bytes.Buffer{}
		if err := tmpl.Execute(buf, data); err != nil {
			return err
		}
		b := buf.Bytes()

		// Run gofmt on the code
		b, err = imports.Process(filename, b, &imports.Options{
			Comments: true,
		})
		if err != nil {
			panic(err)
		}

		// Create the directory if necessary
		dir := filepath.Dir(filename)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0700); err != nil {
				panic(err)
			}
		}

		// Write the file
		if err := ioutil.WriteFile(filename, b, 0644); err != nil {
			panic(err)
		}
	}

	return nil
}
