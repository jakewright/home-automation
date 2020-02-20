package svcdef

import (
	"io/ioutil"
	"os"
)

// FileReader is an interface that wraps ReadFile and SeenFile
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
	SeenFile(filename string) bool
}

type mockFileReader struct {
	files map[string][]byte
	seen  map[string]bool
}

// ReadFile returns the bytes of the file with the given name
func (r *mockFileReader) ReadFile(filename string) ([]byte, error) {
	if r.seen == nil {
		r.seen = map[string]bool{}
	}

	r.seen[filename] = true

	if b, ok := r.files[filename]; ok {
		return b, nil
	}

	return nil, os.ErrNotExist
}

// SeenFile returns true if the file has already been read
func (r *mockFileReader) SeenFile(filename string) bool {
	if r.seen == nil {
		r.seen = map[string]bool{}
	}

	return r.seen[filename]
}

type osFileReader struct {
	seen map[string]bool
}

// ReadFile returns the bytes of the file with the given name
func (r *osFileReader) ReadFile(filename string) ([]byte, error) {
	if r.seen == nil {
		r.seen = map[string]bool{}
	}

	r.seen[filename] = true

	return ioutil.ReadFile(filename)
}

// SeenFile returns true if the file has already been read
func (r *osFileReader) SeenFile(filename string) bool {
	if r.seen == nil {
		r.seen = map[string]bool{}
	}

	return r.seen[filename]
}
