package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	toolsimports "golang.org/x/tools/imports"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/tools/libraries/imports"
)

const pkg = "domain"

var resolver *imports.Resolver

func init() {
	var err error
	resolver, err = imports.NewResolver()
	if err != nil {
		panic(err)
	}
}

type file map[string]*description

type description struct {
	Properties map[string]property `json:"properties"`
}

type property struct {
	Type    string              `json:"type"`
	Min     *float64            `json:"min"`
	Max     *float64            `json:"max"`
	Options []*devicedef.Option `json:"options"`
}

type generator interface {
	Data(im *imports.Manager) (interface{}, error)
	Template() (*template.Template, error)
	Filename() string
}

func generate(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var f file
	if err := json.Unmarshal(b, &f); err != nil {
		panic(err)
	}

	generators := []generator{
		&propertyTypesGenerator{file: f},
	}

	for _, generator := range generators {
		filename := generator.Filename()

		self, err := resolver.Resolve(path, pkg)
		if err != nil {
			panic(err)
		}

		im := imports.NewManager(self)

		// Generate the template data
		data, err := generator.Data(im)
		if err != nil {
			panic(err)
		}

		// If there's nothing to generate, delete the file
		// in case it previously existed.
		if data == nil {
			if err := os.Remove(filename); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				panic(err)
			}

			continue
		}

		// Get the template
		tmpl, err := generator.Template()
		if err != nil {
			panic(err)
		}

		// Generate the code
		buf := &bytes.Buffer{}
		if err := tmpl.Execute(buf, data); err != nil {
			panic(err)
		}
		b := buf.Bytes()

		// Run gofmt on the code
		if true {
			b, err = toolsimports.Process(filename, b, &toolsimports.Options{
				Comments: true,
			})
			if err != nil {
				panic(err)
			}
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
}
