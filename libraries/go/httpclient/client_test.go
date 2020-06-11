package httpclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_methodHelpers(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(fmt.Sprintf("%s %s", r.Method, r.RequestURI)))
		require.NoError(t, err)
	})

	srv := httptest.NewServer(h)
	defer srv.Close()
	c := NewFromHTTPClient(srv.Client(), nil)

	rsp, err := c.Get(context.Background(), srv.URL+"/foo", nil)
	require.NoError(t, err)
	rspBody, err := rsp.BodyString()
	require.NoError(t, err)
	require.Equal(t, "GET /foo", rspBody)

	rsp, err = c.Post(context.Background(), srv.URL+"/foo", nil, nil)
	require.NoError(t, err)
	rspBody, err = rsp.BodyString()
	require.NoError(t, err)
	require.Equal(t, "POST /foo", rspBody)

	rsp, err = c.Put(context.Background(), srv.URL+"/foo", nil, nil)
	require.NoError(t, err)
	rspBody, err = rsp.BodyString()
	require.NoError(t, err)
	require.Equal(t, "PUT /foo", rspBody)

	rsp, err = c.Patch(context.Background(), srv.URL+"/foo", nil, nil)
	require.NoError(t, err)
	rspBody, err = rsp.BodyString()
	require.NoError(t, err)
	require.Equal(t, "PATCH /foo", rspBody)

	rsp, err = c.Delete(context.Background(), srv.URL+"/foo", nil, nil)
	require.NoError(t, err)
	rspBody, err = rsp.BodyString()
	require.NoError(t, err)
	require.Equal(t, "DELETE /foo", rspBody)
}

func TestClient_invalidMethod(t *testing.T) {
	c := New(nil)
	_, err := c.Send(&Request{
		Method: "INVALID",
	}).Response()
	var target InvalidMethodError
	require.True(t, errors.As(err, &target))
}

func TestClient_jsonDecoder(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("ct"))
		body := `{"foo": "bar"}`
		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	})

	srv := httptest.NewServer(h)
	defer srv.Close()
	c := NewFromHTTPClient(srv.Client(), nil)

	// Test base case of not decoding the body
	rsp, err := c.Send(&Request{
		Method: "GET",
		URL:    srv.URL,
	}).Response()
	require.NoError(t, err)
	rspBody, err := rsp.BodyString()
	require.NoError(t, err)
	require.Equal(t, `{"foo": "bar"}`, rspBody)

	// Test inferred decoder from ContentType header
	v := &struct {
		Foo string `json:"foo"`
	}{}
	_, err = c.Send(&Request{
		Method:  "GET",
		URL:     srv.URL,
		Headers: map[string][]string{"ct": {"application/json"}},
	}).DecodeResponse(v)
	require.NoError(t, err)
	require.Equal(t, "bar", v.Foo)

	// Test explicit decoder
	v = &struct {
		Foo string `json:"foo"`
	}{}
	_, err = c.Send(&Request{
		Method:          "GET",
		URL:             srv.URL,
		ResponseDecoder: &DecoderJSON{},
	}).DecodeResponse(v)
	require.NoError(t, err)
	require.Equal(t, "bar", v.Foo)
}
