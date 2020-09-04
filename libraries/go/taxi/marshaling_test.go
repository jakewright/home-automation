package taxi

import (
	"bytes"
	"net/http"
	"testing"

	"gotest.tools/assert"
)

func TestDecodeQuery(t *testing.T) {
	r, err := http.NewRequest("GET", "/baz?foo=bar", nil)
	assert.NilError(t, err)

	var v struct {
		Foo string
	}

	err = DecodeRequest(r, &v)
	assert.NilError(t, err)

	assert.Equal(t, v.Foo, "bar")
}

func TestDecodeBody(t *testing.T) {
	body := []byte("{\"foo\":\"bar\"}")
	r, err := http.NewRequest("POST", "/foo", bytes.NewBuffer(body))
	assert.NilError(t, err)

	var v struct {
		Foo string
	}

	err = DecodeRequest(r, &v)
	assert.NilError(t, err)

	assert.Equal(t, v.Foo, "bar")
}

func TestDecodeIntoMap(t *testing.T) {
	body := []byte("{\"foo\":\"bar\"}")
	r, err := http.NewRequest("GET", "/baz?baz=qux", bytes.NewBuffer(body))
	assert.NilError(t, err)

	var v map[string]string

	err = DecodeRequest(r, &v)
	assert.NilError(t, err)

	assert.Equal(t, v["foo"], "bar")
	assert.Equal(t, v["baz"], "qux")
}

func TestDecodeComplexParamNames(t *testing.T) {
	body := []byte("{\"animal_color\":\"black\"}")
	r, err := http.NewRequest("GET", "/foo?house_name=Buckingham%20Palace", bytes.NewBuffer(body))
	assert.NilError(t, err)

	var v struct {
		AnimalColor string `json:"animal_color"`
		HouseName   string `json:"house_name"`
	}

	err = DecodeRequest(r, &v)
	assert.NilError(t, err)

	assert.Equal(t, v.AnimalColor, "black")
	assert.Equal(t, v.HouseName, "Buckingham Palace")
}
