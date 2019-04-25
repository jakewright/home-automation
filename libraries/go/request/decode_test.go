package request

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/assert"
)

func TestDecodeMux(t *testing.T) {
	r, err := http.NewRequest("GET", "/foo", nil)
	assert.NilError(t, err)

	r = mux.SetURLVars(r, map[string]string{"foo": "bar"})

	var v struct {
		Foo string
	}

	err = Decode(r, &v)
	assert.NilError(t, err)

	assert.Equal(t, v.Foo, "bar")
}

func TestDecodeQuery(t *testing.T) {
	r, err := http.NewRequest("GET", "/baz?foo=bar", nil)
	assert.NilError(t, err)

	var v struct {
		Foo string
	}

	err = Decode(r, &v)
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

	err = Decode(r, &v)
	assert.NilError(t, err)

	assert.Equal(t, v.Foo, "bar")
}

func TestDecodeIntoMap(t *testing.T) {
	body := []byte("{\"foo\":\"bar\"}")
	r, err := http.NewRequest("GET", "/baz?baz=qux", bytes.NewBuffer(body))
	assert.NilError(t, err)

	r = mux.SetURLVars(r, map[string]string{"quz": "cog"})

	var v map[string]string

	err = Decode(r, &v)
	assert.NilError(t, err)

	assert.Equal(t, v["foo"], "bar")
	assert.Equal(t, v["baz"], "qux")
	assert.Equal(t, v["quz"], "cog")
}

func TestDecodeComplexParamNames(t *testing.T) {
	body := []byte("{\"animal_color\":\"black\"}")
	r, err := http.NewRequest("GET", "/foo?house_name=Buckingham%20Palace", bytes.NewBuffer(body))
	assert.NilError(t, err)

	r = mux.SetURLVars(r, map[string]string{"favorite_number": "3"})

	var v struct {
		AnimalColor    string `json:"animal_color"`
		HouseName      string `json:"house_name"`
		FavoriteNumber int    `json:"favorite_number"`
	}

	err = Decode(r, &v)
	assert.NilError(t, err)

	assert.Equal(t, v.AnimalColor, "black")
	assert.Equal(t, v.HouseName, "Buckingham Palace")
	assert.Equal(t, v.FavoriteNumber, 3)
}
