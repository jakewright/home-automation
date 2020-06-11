package httpclient

import (
	"context"
	"net/http"
	"time"
)

// DefaultTimeout is the default time limit for requests made by the client.
const DefaultTimeout = 30 * time.Second

// DefaultStatusValidator returns true for 2xx statuses, otherwise false.
var DefaultStatusValidator = func(status int) bool {
	return status >= 200 && status < 300
}

// VoidStatusValidator always returns true.
var VoidStatusValidator = func(int) bool { return true }

// Client is a wrapper around the standard library HTTP client
type Client struct {
	validateStatus  func(int) bool
	requestEncoder  Encoder
	responseDecoder Decoder
	baseClient      httpClient
}

// Options allows client defaults to be overridden
type Options struct {
	// ValidateStatus is a function that validates the status code of the HTTP
	// response. If the function returns false, Do() will return an error.
	// If nil, DefaultStatusValidator is used.
	ValidateStatus func(int) bool

	// RequestEncoder is the Encoding used to marshal request bodies. This can
	// be overridden per-request. If nil, EncodingJSON is used.
	RequestEncoder Encoder

	// ResponseDecoder is the Encoding used to unmarshal response bodies. This
	// can be overridden per-request. If nil, the Encoding is inferred from
	// the response's ContentType header.
	ResponseDecoder Decoder

	// Timeout specifies the time limit for requests made by the client.
	// If nil, DefaultTimeout is used.
	// This option is ignored if NewFromHTTPClient is used.
	Timeout time.Duration
}

// New returns a new Client.
func New(opts *Options) *Client {
	if opts == nil {
		opts = &Options{}
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	baseClient := &http.Client{
		Timeout: timeout,
	}

	return NewFromHTTPClient(baseClient, opts)
}

// NewFromHTTPClient returns a new Client that wraps baseClient.
func NewFromHTTPClient(baseClient httpClient, opts *Options) *Client {
	if opts == nil {
		opts = &Options{}
	}

	return &Client{
		validateStatus:  opts.ValidateStatus,
		requestEncoder:  opts.RequestEncoder,
		responseDecoder: opts.ResponseDecoder,
		baseClient:      baseClient,
	}
}

// httpClient is an interface that http.Client{} implements
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Send performs the HTTP request and returns a Future
func (c *Client) Send(request *Request) *Future {
	done := make(chan struct{})
	ftr := &Future{done: done, request: request}

	go func() {
		defer close(done)
		ftr.response, ftr.err = c.do(request)
	}()

	return ftr
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, url string, v interface{}) (*Response, error) {
	r := &Request{Ctx: ctx, Method: http.MethodGet, URL: url}
	return c.Send(r).DecodeResponse(v)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, url string, body interface{}, v interface{}) (*Response, error) {
	r := &Request{Ctx: ctx, Method: http.MethodPost, URL: url, Body: body}
	return c.Send(r).DecodeResponse(v)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, url string, body interface{}, v interface{}) (*Response, error) {
	r := &Request{Ctx: ctx, Method: http.MethodPut, URL: url, Body: body}
	return c.Send(r).DecodeResponse(v)
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, url string, body interface{}, v interface{}) (*Response, error) {
	r := &Request{Ctx: ctx, Method: http.MethodPatch, URL: url, Body: body}
	return c.Send(r).DecodeResponse(v)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, url string, body interface{}, v interface{}) (*Response, error) {
	r := &Request{Ctx: ctx, Method: http.MethodDelete, URL: url, Body: body}
	return c.Send(r).DecodeResponse(v)
}

func (c *Client) do(request *Request) (*Response, error) {
	if err := request.validate(); err != nil {
		return nil, err
	}

	//  Build the HTTP request
	req, err := request.build(c.encoder())
	if err != nil {
		return nil, err
	}

	// Make the request
	rsp, err := c.baseClient.Do(req)
	if err != nil {
		return nil, err
	}

	// From this point on, all return values should return response, even if there's an error
	// so that the caller can see all of the information about the response.
	response := &Response{Response: rsp}

	// Validate the status
	if !c.statusValidator()(rsp.StatusCode) {
		return response, BadStatusError(rsp.StatusCode)
	}

	return response, nil
}

func (c *Client) statusValidator() func(int) bool {
	if c.validateStatus != nil {
		return c.validateStatus
	}

	return DefaultStatusValidator
}

func (c *Client) encoder() Encoder {
	if c.requestEncoder != nil {
		return c.requestEncoder
	}

	return EncoderJSON{}
}
