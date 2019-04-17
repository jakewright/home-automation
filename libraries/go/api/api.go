package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Requester is an interface for making HTTP requests
type Requester interface {
	Do(request *Request, resultData interface{}) (*Response, error)
	Get(url string, response interface{}) (*Response, error)
	Put(url string, body map[string]interface{}, response interface{}) (*Response, error)
	Patch(url string, body map[string]interface{}, response interface{}) (*Response, error)
}

type Client struct {
	// Base is the base URL used for relative requests
	Base string

	// ValidateStatus is a function that validates the status code of the HTTP response
	ValidateStatus func(int) bool

	// Envelope is the name of the data field in the response
	Envelope string
}

var DefaultClient Requester

func mustGetDefaultClient() Requester {
	if DefaultClient == nil {
		panic("Default HTTP client used before being set")
	}

	return DefaultClient
}

func Get(url string, response interface{}) (*Response, error) {
	return mustGetDefaultClient().Get(url, response)
}
func Put(url string, body map[string]interface{}, response interface{}) (*Response, error) {
	return mustGetDefaultClient().Put(url, body, response)
}
func Patch(url string, body map[string]interface{}, response interface{}) (*Response, error) {
	return mustGetDefaultClient().Patch(url, body, response)
}

// New returns a new API Client
func New(base string, envelope string) (Requester, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	return &Client{
		Base: strings.TrimRight(u.String(), "/"),
		ValidateStatus: func(status int) bool {
			return status >= 200 && status < 300
		},
		Envelope: envelope,
	}, nil
}

// Request holds the information needed to make an HTTP request
type Request struct {
	Method string
	URL    string
	Body   map[string]interface{}
}

type Response struct {
	*http.Response
	Body string
}

// Do performs the HTTP request. Relative URLs will have the base URL prepended. An error
// will be thrown if the response does not have a JSON content-type or if the status code is
// not valid. The entire response will be returned but the JSON will be unmarshalled into the second argument.
func (c Client) Do(request *Request, v interface{}) (*Response, error) {
	if !validMethod(request.Method) {
		return nil, fmt.Errorf("method %q is not a valid HTTP method", request.Method)
	}

	// Convert the body to a byte array
	var jsonBody []byte
	var err error
	if request.Body != nil {
		fmt.Println("Body is not nil")
		jsonBody, err = json.Marshal(request.Body)
		if err != nil {
			return nil, err
		}
	}

	// If the URL is relative, prepend the base URL
	absURL := request.URL
	u, err := url.Parse(request.URL)
	if !u.IsAbs() {
		absURL = fmt.Sprintf("%s/%s", c.Base, strings.TrimLeft(request.URL, "/"))
	}

	client := &http.Client{}
	req, err := http.NewRequest(request.Method, absURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	// Make the request
	rawRsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// From this point on, all return values should return rsp, even if there's an error
	// so that the caller can see all of the information about the response.
	rsp := &Response{Response: rawRsp}

	// Read the body into a byte array
	defer rawRsp.Body.Close()
	body, err := ioutil.ReadAll(rawRsp.Body)
	if err != nil {
		return rsp, err
	}
	rsp.Body = string(body)

	// Validate the status
	if !c.ValidateStatus(rawRsp.StatusCode) {
		return rsp, fmt.Errorf("%s %s\n"+
			"request failed with status %s\n"+
			"%s", request.Method, absURL, rawRsp.Status, body)
	}

	// Validate the content type
	contentType := rawRsp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "json") {
		return rsp, fmt.Errorf("content type %s not supported", contentType)
	}

	if c.Envelope != "" {
		// Given we know it's enveloped, we can unmarshal into a map
		var rspWrapper map[string]interface{}
		if err = json.Unmarshal(body, &rspWrapper); err != nil {
			return rsp, err
		}

		// Extract the inner field and marshal back into JSON
		innerBytes, err := json.Marshal(rspWrapper[c.Envelope])
		if err != nil {
			return rsp, err
		}

		// Allow the unmarshal function to deal with reflection and set the value of &response
		if err = json.Unmarshal(innerBytes, &v); err != nil {
			return rsp, err
		}
	} else {
		if err = json.Unmarshal(body, &v); err != nil {
			return rsp, err
		}
	}

	return &Response{Response: rawRsp, Body: string(body)}, nil
}

// Get performs a GET request
func (c Client) Get(url string, response interface{}) (*Response, error) {
	r := Request{Method: http.MethodGet, URL: url}
	return c.Do(&r, response)
}

// Put performs a PUT request
func (c Client) Put(url string, body map[string]interface{}, response interface{}) (*Response, error) {
	r := Request{Method: http.MethodPut, URL: url, Body: body}
	return c.Do(&r, response)
}

// Patch performs a PATCH request
func (c Client) Patch(url string, body map[string]interface{}, response interface{}) (*Response, error) {
	r := Request{Method: http.MethodPatch, URL: url, Body: body}
	return c.Do(&r, response)
}

func validMethod(method string) bool {
	switch method {
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodPatch:
	case http.MethodDelete:
	case http.MethodConnect:
	case http.MethodOptions:
	case http.MethodTrace:
	default:
		return false
	}

	return true
}
