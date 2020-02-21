package svcdef

import (
	"testing"

	"gotest.tools/assert"
)

func TestParser_Parse_emptyFile(t *testing.T) {
	fr := &mockFileReader{
		files: map[string][]byte{
			"test.def": []byte(``),
		},
	}

	f, err := NewParser(fr).Parse("test.def")
	assert.NilError(t, err)

	expected := &File{
		Path: "test.def",
	}

	assert.DeepEqual(t, expected, f)
}

func TestParser_Parse_simpleRPC(t *testing.T) {
	input := []byte(`
foo = "bar"

service Test {
	path = "service.test"

	rpc GetFoo(GetFooRequest) GetFooResponse {
		method = "GET"
		path = "/foo"
	}
}

message GetFooRequest {
}

message GetFooResponse {
	string bar
}
`)

	fr := &mockFileReader{
		files: map[string][]byte{
			"test.def": input,
		},
	}

	f, err := NewParser(fr).Parse("test.def")
	assert.NilError(t, err)

	messages := []*Message{
		{
			Name:          "GetFooRequest",
			QualifiedName: ".GetFooRequest",
		},
		{
			Name:          "GetFooResponse",
			QualifiedName: ".GetFooResponse",
			Fields: []*Field{
				{
					Name: "bar",
					Type: &Type{
						Name:     "string",
						Original: "string",
					},
				},
			},
		},
	}

	expected := &File{
		Path: "test.def",
		Service: &Service{
			Name: "Test",
			RPCs: []*RPC{
				{
					Name: "GetFoo",
					InputType: &Type{
						Name:      "GetFooRequest",
						Original:  "GetFooRequest",
						Qualified: ".GetFooRequest",
					},
					OutputType: &Type{
						Name:      "GetFooResponse",
						Original:  "GetFooResponse",
						Qualified: ".GetFooResponse",
					},
					Options: map[string]interface{}{
						"method": "GET",
						"path":   "/foo",
					},
				},
			},
			Options: map[string]interface{}{
				"path": "service.test",
			},
		},
		Messages:     messages,
		FlatMessages: messages,
		Options: map[string]interface{}{
			"foo": "bar",
		},
	}

	assert.DeepEqual(t, expected, f)
}

func TestParser_Parse_nestedMessage(t *testing.T) {
	input := []byte(`message Foo {
	message Bar {
	}

	Bar bar
}`)

	fr := &mockFileReader{
		files: map[string][]byte{
			"test.def": input,
		},
	}

	actual, err := NewParser(fr).Parse("test.def")
	assert.NilError(t, err)

	bar := &Message{
		Name:          "Bar",
		QualifiedName: ".Foo.Bar",
	}

	foo := &Message{
		Name:          "Foo",
		QualifiedName: ".Foo",
		Fields: []*Field{
			{
				Name: "bar",
				Type: &Type{
					Name:      "Bar",
					Original:  "Bar",
					Qualified: ".Foo.Bar",
				},
			},
		},
		Nested: []*Message{bar},
	}

	expected := &File{
		Path:         "test.def",
		Messages:     []*Message{foo},
		FlatMessages: []*Message{foo, bar},
	}

	assert.DeepEqual(t, expected, actual)
}

func TestParser_Parse_fieldOptions(t *testing.T) {
	input := []byte(`message Foo {
	*[]string bar (required,foo="bar",bat=5,baz=false)
}`)

	fr := &mockFileReader{
		files: map[string][]byte{
			"test.def": input,
		},
	}

	actual, err := NewParser(fr).Parse("test.def")
	assert.NilError(t, err)

	messages := []*Message{
		{
			Name:          "Foo",
			QualifiedName: ".Foo",
			Fields: []*Field{
				{
					Name: "bar",
					Type: &Type{
						Name:     "string",
						Original: "*[]string",
						Repeated: true,
						Optional: true,
					},
					Options: map[string]interface{}{
						"required": true,
						"foo":      "bar",
						"bat":      int64(5),
						"baz":      false,
					},
				},
			},
		},
	}

	expected := &File{
		Path:         "test.def",
		Messages:     messages,
		FlatMessages: messages,
	}

	assert.DeepEqual(t, expected, actual)
}

func TestParser_Parse_mapType(t *testing.T) {
	input := []byte(`message Foo {
	message Bar {
	}

	map[Baz]string foo
	*[]map[map[string]int][]Bar bar
}

message Baz {
}

`)

	fr := &mockFileReader{
		files: map[string][]byte{
			"test.def": input,
		},
	}

	actual, err := NewParser(fr).Parse("test.def")
	assert.NilError(t, err)

	bar := &Message{
		Name:          "Bar",
		QualifiedName: ".Foo.Bar",
	}

	foo := &Message{
		Name:          "Foo",
		QualifiedName: ".Foo",
		Fields: []*Field{
			{
				Name: "foo",
				Type: &Type{
					Name:     "map",
					Original: "map[Baz]string",
					Map:      true,
					MapKey: &Type{
						Name:      "Baz",
						Original:  "Baz",
						Qualified: ".Baz",
					},
					MapValue: &Type{
						Name:     "string",
						Original: "string",
					},
				},
			},
			{
				Name: "bar",
				Type: &Type{
					Name:     "map",
					Original: "*[]map[map[string]int][]Bar",
					Map:      true,
					MapKey: &Type{
						Name:     "map",
						Original: "map[string]int",
						Map:      true,
						MapKey: &Type{
							Name:     "string",
							Original: "string",
						},
						MapValue: &Type{
							Name:     "int",
							Original: "int",
						},
					},
					MapValue: &Type{
						Name:      "Bar",
						Original:  "[]Bar",
						Qualified: ".Foo.Bar",
						Repeated:  true,
					},
					Repeated: true,
					Optional: true,
				},
			},
		},
		Nested: []*Message{bar},
	}

	baz := &Message{
		Name:          "Baz",
		QualifiedName: ".Baz",
	}

	expected := &File{
		Path:         "test.def",
		Messages:     []*Message{foo, baz},
		FlatMessages: []*Message{foo, bar, baz},
	}

	assert.DeepEqual(t, expected, actual)
}

func TestParser_Parse_importedType(t *testing.T) {
	file1 := []byte(`import bar "../service.bar/bar.def"

service Svc {
	rpc GetFoo(bar.Foo) bar.Bar {
	}
}

message Msg {
	bar.Foo foo
	bar.Bar bar
}
`)

	file2 := []byte(`message Foo {
}

message Bar {
}`)

	fr := &mockFileReader{
		files: map[string][]byte{
			"service.foo/foo.def": file1,
			"service.bar/bar.def": file2,
		},
	}

	f, err := NewParser(fr).Parse("service.foo/foo.def")
	assert.NilError(t, err)

	foo := &Message{Name: "Foo", QualifiedName: "bar.Foo"}
	bar := &Message{Name: "Bar", QualifiedName: "bar.Bar"}
	msg := &Message{
		Name:          "Msg",
		QualifiedName: ".Msg",
		Fields: []*Field{
			{
				Name: "foo",
				Type: &Type{Name: "bar.Foo", Original: "bar.Foo", Qualified: "bar.Foo"},
			},
			{
				Name: "bar",
				Type: &Type{Name: "bar.Bar", Original: "bar.Bar", Qualified: "bar.Bar"},
			},
		},
	}

	expected := &File{
		Path: "service.foo/foo.def",
		Imports: map[string]*Import{
			"bar": {
				File: &File{
					Path:         "service.bar/bar.def",
					Messages:     []*Message{foo, bar},
					FlatMessages: []*Message{foo, bar},
				},
				Alias: "bar",
				Path:  "../service.bar/bar.def",
			},
		},
		Service: &Service{
			Name: "Svc",
			RPCs: []*RPC{
				{
					Name:       "GetFoo",
					InputType:  &Type{Name: "bar.Foo", Original: "bar.Foo", Qualified: "bar.Foo"},
					OutputType: &Type{Name: "bar.Bar", Original: "bar.Bar", Qualified: "bar.Bar"},
				},
			},
		},
		Messages:     []*Message{msg},
		FlatMessages: []*Message{msg, foo, bar},
	}

	assert.DeepEqual(t, expected, f)
}

func TestParser_Parse_circularImport(t *testing.T) {
	fr := &mockFileReader{
		files: map[string][]byte{
			"file1.def": []byte(`import file2 "file2.def"`),
			"file2.def": []byte(`import file3 "file3.def"`),
			"file3.def": []byte(`import file1 "file1.def"`),
		},
	}

	_, err := NewParser(fr).Parse("file1.def")
	assert.Error(t, err, "file1.def -> file2.def -> file3.def -> file1.def: circular import")
}
