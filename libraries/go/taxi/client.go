package taxi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jakewright/patch"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

const requestTimeout = 10 * time.Second

var (
	defaultDispatcher Dispatcher
	once              = &sync.Once{}
)

// Dispatcher is the interface for making remote procedure calls
type Dispatcher interface {
	Dispatch(*RPC) *Future
	Get(ctx context.Context, url string, body interface{}, v interface{}) error
	Post(ctx context.Context, url string, body interface{}, v interface{}) error
	Put(ctx context.Context, url string, body interface{}, v interface{}) error
	Patch(ctx context.Context, url string, body interface{}, v interface{}) error
	Delete(ctx context.Context, url string, body interface{}, v interface{}) error
}

// RPC represents a remote procedure call
type RPC struct {
	Method string
	URL    string
	Body   interface{}
	ctx    context.Context
}

// NewRPC returns a new RPC given a context, method, url and body. The body
// will be marshaled into JSON by the client when dispatching the RPC.
func NewRPC(ctx context.Context, method, url string, body interface{}) (*RPC, error) {
	if !validMethod(method) {
		return nil, fmt.Errorf("taxi: invalid method: %q", method)
	}

	if ctx == nil {
		return nil, errors.New("taxi: nil context")
	}

	return &RPC{
		Method: method,
		URL:    url,
		Body:   body,
		ctx:    ctx,
	}, nil
}

// SetDefaultDispatcher sets the default client once and only once
func SetDefaultDispatcher(d Dispatcher) {
	once.Do(func() { defaultDispatcher = d })
}

// mustGetDefaultDispatcher returns the default dispatcher and panics if it is nil
func mustGetDefaultDispatcher() Dispatcher {
	if defaultDispatcher == nil {
		panic(fmt.Errorf("no default Dispatcher set"))
	}

	return defaultDispatcher
}

// Dispatch makes a request using the default client
func Dispatch(rpc *RPC) *Future {
	return mustGetDefaultDispatcher().Dispatch(rpc)
}

// Get dispatches a GET RPC using the default dispatcher
func Get(ctx context.Context, url string, body interface{}, v interface{}) error {
	return mustGetDefaultDispatcher().Get(ctx, url, body, v)
}

// Post dispatches a POST RPC using the default dispatcher
func Post(ctx context.Context, url string, body interface{}, v interface{}) error {
	return mustGetDefaultDispatcher().Post(ctx, url, body, v)
}

// Put dispatches a PUT RPC using the default dispatcher
func Put(ctx context.Context, url string, body interface{}, v interface{}) error {
	return mustGetDefaultDispatcher().Put(ctx, url, body, v)
}

// Patch dispatches a PATCH RPC using the default dispatcher
func Patch(ctx context.Context, url string, body interface{}, v interface{}) error {
	return mustGetDefaultDispatcher().Patch(ctx, url, body, v)
}

// Delete dispatches a DELETE RPC using the default dispatcher
func Delete(ctx context.Context, url string, body interface{}, v interface{}) error {
	return mustGetDefaultDispatcher().Delete(ctx, url, body, v)
}

// Doer is an interface that http.Client implements
type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}

// Client dispatches RPC requests
type Client struct {
	p *patch.Client
}

// NewClient returns an initialised Client
func NewClient() *Client {
	return NewClientFromDoer(&http.Client{
		Timeout: requestTimeout,
	})
}

// NewClientFromDoer returns an initialised client, using the given Doer
// to make the HTTP requests. The standard http.Client implements the
// Doer interface.
func NewClientFromDoer(doer Doer) *Client {
	return &Client{
		p: patch.NewFromBaseClient(doer),
	}
}

// Dispatch makes a request and returns a Future
// that represents the in-flight request
func (c *Client) Dispatch(rpc *RPC) *Future {
	req := &patch.Request{
		Ctx:    rpc.ctx,
		Method: rpc.Method,
		URL:    rpc.URL,
		Body:   rpc.Body,
	}

	done := make(chan struct{})
	ftr := &Future{done: done}

	go func() {
		defer close(done)

		errParams := map[string]string{
			"method": rpc.Method,
			"url":    rpc.URL,
		}

		// Block until the response is returned
		rsp, err := c.p.Send(req).Response()

		ftr.decode = func(v interface{}) error {
			if err != nil {
				return err
			}

			var wrapper struct {
				Error string          `json:"error"`
				Data  json.RawMessage `json:"data"`
			}

			if err := rsp.DecodeJSON(&wrapper); err != nil {
				return oops.WithMessage(err, "received %d %s from server but could not decode response", rsp.StatusCode, rsp.Status, errParams)
			}

			if wrapper.Error != "" {
				return oops.FromHTTPStatus(rsp.StatusCode, wrapper.Error, errParams)
			}

			if v == nil {
				return nil
			}

			if err := json.Unmarshal(wrapper.Data, v); err != nil {
				return oops.WithMetadata(err, errParams)
			}

			return nil
		}
	}()

	return ftr
}

// Get dispatches a GET RPC
func (c *Client) Get(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodGet, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}

// Post dispatches a POST RPC
func (c *Client) Post(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPost, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}

// Put dispatches a PUT RPC
func (c *Client) Put(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPut, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}

// Patch dispatches a PATCH RPC
func (c *Client) Patch(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodPatch, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
}

// Delete dispatches a DELETE RPC
func (c *Client) Delete(ctx context.Context, url string, body interface{}, v interface{}) error {
	r := &RPC{Method: http.MethodDelete, URL: url, Body: body, ctx: ctx}
	return c.Dispatch(r).DecodeResponse(v)
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

func validMethod(method string) bool {
	switch method {
	case http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete:
		return true
	}

	return false
}
