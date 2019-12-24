package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/jakewright/home-automation/tools/protoc-gen-jrpc/gen"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	req := readRequest(os.Stdin)
	g := gen.Generator{}
	rsp := g.Generate(req)
	writeResponse(os.Stdout, rsp)
}

func readRequest(r io.Reader) *plugin.CodeGeneratorRequest {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	var req plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(b, &req); err != nil {
		panic(err)
	}

	if len(req.FileToGenerate) == 0 {
		panic("No files to generate")
	}

	return &req
}

func writeResponse(w io.Writer, rsp *plugin.CodeGeneratorResponse) {
	b, err := proto.Marshal(rsp)
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		panic(err)
	}
}
