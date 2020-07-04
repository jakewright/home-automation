package taxi

import (
	"context"
	"net/http"
	"time"
)

const requestTimeout = 10 * time.Second

// Dispatcher is the interface for making remote procedure calls
type Dispatcher interface {
	Dispatch(ctx context.Context, rpc *RPC) *Future
	Get(ctx context.Context, url string, body interface{}, v interface{}) error
	Post(ctx context.Context, url string, body interface{}, v interface{}) error
	Put(ctx context.Context, url string, body interface{}, v interface{}) error
	Patch(ctx context.Context, url string, body interface{}, v interface{}) error
	Delete(ctx context.Context, url string, body interface{}, v interface{}) error
}

// Doer is an interface that http.Client implements
type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}

// Client dispatches RPC requests
type Client struct {
	base Doer
}

// NewClient returns an initialised Client
func NewClient() *Client {
	return NewClientUsing(&http.Client{
		Timeout: requestTimeout,
	})
}

// NewClientUsing returns an initialised client, using the given Doer
// to make the HTTP requests. The standard http.Client implements the
// Doer interface.
func NewClientUsing(doer Doer) *Client {
	return &Client{
		base: doer,
	}
}

// Dispatch makes a request and returns a Future
// that represents the in-flight request
func (c *Client) Dispatch(ctx context.Context, rpc *RPC) *Future {
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

// Get dispatches a GET RPC
func (c *Client) Get(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodGet, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

// Post dispatches a POST RPC
func (c *Client) Post(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPost, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

// Put dispatches a PUT RPC
func (c *Client) Put(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPut, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

// Patch dispatches a PATCH RPC
func (c *Client) Patch(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPatch, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

// Delete dispatches a DELETE RPC
func (c *Client) Delete(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodDelete, URL: url, Body: body}
	return c.Dispatch(ctx, r).DecodeResponse(v)
}

func (c *Client) do(ctx context.Context, rpc *RPC) (*http.Response, error) {
	req, err := rpc.ToRequest(ctx)
	if err != nil {
		return nil, err
	}

	return c.base.Do(req)
}

// Future represents an in-flight remote procedure call
type Future struct {
	done   <-chan struct{}
	decode func(v interface{}) error
}

// Wait will block until the response has been received
func (f *Future) Wait() error {
	return f.DecodeResponse(nil)
}

// DecodeResponse will block until the response has been received, and then
// decode the JSON body into the given argument.
func (f *Future) DecodeResponse(v interface{}) error {
	<-f.done
	return f.decode(v)
}
