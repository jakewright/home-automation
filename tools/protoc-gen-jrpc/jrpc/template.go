package jrpc

import "text/template"

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.New("jrpc").Parse(templateText)
	if err != nil {
		panic(err)
	}
}

const templateText = `// Code generated by protoc-gen-jrpc. DO NOT EDIT.

{{ range .PackageComment }}
// {{ . }}
{{ end }}
package {{ .PackageName }}

import (
{{ range .Imports }}
    {{ .Alias }} "{{ .Path }}"
{{ end }}
)

type Builder struct {
{{ range .Service.RPCs }}
    {{ .Name }} func(request *{{ .RequestType }}) (*{{ .ResponseType }}, error)
{{ end }}
}
`