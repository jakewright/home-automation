package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	req := readRequest(os.Stdin)

	rsp, err := generate(req)
	if err != nil {
		panic(err)
	}
	if rsp == nil {
		return // Nothing to do
	}

	writeResponse(os.Stdout, rsp)
}

func readRequest(r io.Reader) *plugin_go.CodeGeneratorRequest {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	var req plugin_go.CodeGeneratorRequest
	if err := proto.Unmarshal(b, &req); err != nil {
		panic(err)
	}

	if len(req.FileToGenerate) == 0 {
		panic("No files to generate")
	}

	return &req
}

func writeResponse(w io.Writer, rsp *plugin_go.CodeGeneratorResponse) {
	b, err := proto.Marshal(rsp)
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		panic(err)
	}
}
