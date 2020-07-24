package taxi

import (
	"context"
	"net/http"
	"net/http/httptest"
)

// MockClient dispatches RPCs by passing the request directly to an internal
// http.Handler. It does not make network requests. This is useful in unit
// tests as it allows endpoints to be mocked.
type MockClient struct {
	Handler http.Handler
}

var _ Dispatcher = (*MockClient)(nil)

// Dispatch converts the RPC into an http request and gives it to the client's
// handler to handle. It returns a Future which will resolve to the response
// given by the handler.
func (c *MockClient) Dispatch(ctx context.Context, rpc *RPC) *Future {
	done := make(chan struct{})
	ftr := &Future{done: done}

	go func() {
		defer close(done)
		rsp, err := c.do(ctx, rpc)

		ftr.decode = func(v interface{}) error {
			if err != nil {
				return err
			}

			return decodeResponse(rsp, v)
		}
	}()

	return ftr
}

func (c *MockClient) do(ctx context.Context, rpc *RPC) (*http.Response, error) {
	req, err := rpc.ToRequest(ctx)
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	c.Handler.ServeHTTP(w, req)

	return w.Result(), nil
}

// Get dispatches a GET RPC
func (c *MockClient) Get(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodGet, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

// Post dispatches a POST RPC
func (c *MockClient) Post(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPost, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

// Put dispatches a PUT RPC
func (c *MockClient) Put(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPut, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

// Patch dispatches a PATCH RPC
func (c *MockClient) Patch(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPatch, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

// Delete dispatches a DELETE RPC
func (c *MockClient) Delete(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodDelete, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}
