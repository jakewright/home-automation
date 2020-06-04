package rpc

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// Client is an interface for making RPCs
type Client interface {
	Do(ctx context.Context, request *Request, resultData interface{}) (*Response, error)
}

var (
	defaultClient Client
	once          = &sync.Once{}
)

// SetDefaultClient sets the default client once and only once
func SetDefaultClient(c Client) {
	once.Do(func() { defaultClient = c })
}

func mustGetDefaultClient(req *Request) Client {
	if defaultClient == nil {
		panic(fmt.Errorf("no default RPC client set for request to %s", req.URL))
	}

	return defaultClient
}

// Do makes a request using the default client
func Do(ctx context.Context, r *Request, response interface{}) (*Response, error) {
	return mustGetDefaultClient(r).Do(ctx, r, response)
}

// Get makes GET requests using the default client
func Get(ctx context.Context, url string, response interface{}) (*Response, error) {
	r := &Request{Method: http.MethodGet, URL: url}
	return mustGetDefaultClient(r).Do(ctx, r, response)
}

// Put makes PUT requests using the default client
func Put(ctx context.Context, url string, body map[string]interface{}, response interface{}) (*Response, error) {
	r := &Request{Method: http.MethodPut, URL: url, Body: body}
	return mustGetDefaultClient(r).Do(ctx, r, response)
}

// Patch makes PATCH requests using the default client
func Patch(ctx context.Context, url string, body map[string]interface{}, response interface{}) (*Response, error) {
	r := &Request{Method: http.MethodPatch, URL: url, Body: body}
	return mustGetDefaultClient(r).Do(ctx, r, response)
}

// Request holds the information needed to make an HTTP request
type Request struct {
	Method string
	URL    string
	Body   interface{}
}

// Response wraps the http.Response returned from the request
type Response struct {
	*http.Response
	Body []byte
}
