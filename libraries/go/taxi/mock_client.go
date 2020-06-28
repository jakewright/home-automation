package taxi

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

// MockClient dispatches RPCs by passing the request directly to an internal
// http.Handler. It does not make network requests. This is useful in unit
// tests as it allows endpoints to be mocked.
type MockClient struct {
	handler http.Handler
}

// NewMockClient returns a new mock client using the given handler
func NewMockClient(handler http.Handler) *MockClient {
	return &MockClient{handler}
}

// Dispatch converts the RPC into an http request and gives it to the client's
// handler to handle. It returns a Future which will resolve to the response
// given by the handler.
func (c *MockClient) Dispatch(rpc *RPC) *Future {
	done := make(chan struct{})
	ftr := &Future{
		done: done,
	}

	go func() {
		defer close(done)

		errParams := map[string]string{
			"method": rpc.Method,
			"url":    rpc.URL,
		}

		rsp, err := c.do(rpc)

		ftr.decode = func(v interface{}) error {
			if err != nil {
				return oops.WithMetadata(err, errParams)
			}

			body, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return oops.WithMessage(err, "failed to read response body", errParams)
			}

			var wrapper struct {
				Error string          `json:"error"`
				Data  json.RawMessage `json:"data"`
			}

			if err := json.Unmarshal(body, &wrapper); err != nil {
				return oops.WithMessage(err, "failed to unmarshal body", errParams)
			}

			if wrapper.Error != "" {
				return oops.FromHTTPStatus(rsp.StatusCode, wrapper.Error, errParams)
			}

			if err := json.Unmarshal(wrapper.Data, v); err != nil {
				return oops.WithMetadata(err, errParams)
			}

			return nil
		}
	}()

	return ftr
}

func (c *MockClient) do(rpc *RPC) (*http.Response, error) {
	bodyBytes, err := json.Marshal(rpc.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(rpc.Method, rpc.URL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	c.handler.ServeHTTP(w, req)

	return w.Result(), nil
}

// Get dispatches a GET RPC
func (c *MockClient) Get(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodGet, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}

// Post dispatches a POST RPC
func (c *MockClient) Post(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPost, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}

// Put dispatches a PUT RPC
func (c *MockClient) Put(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPut, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}

// Patch dispatches a PATCH RPC
func (c *MockClient) Patch(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPatch, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}

// Delete dispatches a DELETE RPC
func (c *MockClient) Delete(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodDelete, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}
