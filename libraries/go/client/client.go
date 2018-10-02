package client

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
	Request(request *Request, resultData interface{}) (*http.Response, error)
	Get(url string, resultData interface{}) (*http.Response, error)
	Patch(url string, body map[string]interface{}, resultData interface{}) (*http.Response, error)
}

type apiClient struct {
	base *url.URL
}

// New returns a new API Client
func New(base string) (Requester, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	return &apiClient{
		base: u,
	}, nil
}

// Request holds the information needed to make an HTTP request
type Request struct {
	Method string
	URL    string
	Body   map[string]interface{}
}

// Request performs the HTTP request. Relative URLs will have the base URL prepended. An error
// will be thrown if the response does not have a JSON content-type or if the status code is
// not in the 200 range. The entire response will be returned but the data field of the response
// will be unmarshalled into the second argument.
func (c apiClient) Request(request *Request, rspData interface{}) (*http.Response, error) {
	// Convert the body to a byte array
	jsonBody, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}

	// If the URL is relative, prepend the base URL
	u, err := url.Parse(request.URL)
	resolvedURL := c.base.ResolveReference(u)

	fmt.Println(resolvedURL.String())

	client := &http.Client{}
	req, err := http.NewRequest(request.Method, resolvedURL.String(), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	// Make the request
	rawRsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Validate the content type
	contentType := rawRsp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "json") {
		return nil, fmt.Errorf("content type %s not supported", contentType)
	}

	// Read the body into a byte array
	defer rawRsp.Body.Close()
	body, err := ioutil.ReadAll(rawRsp.Body)
	if err != nil {
		return nil, err
	}

	// The response will be enveloped in a data field so unmarshal
	// it into that format and then discard the outer layer.
	rsp := struct {
		Data interface{}
	}{
		Data: rspData,
	}

	// Decode the byte array into the rsp struct. The caller will maintain
	// access to rspData because, as an interface, it's passed by reference.
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return nil, err
	}

	return rawRsp, nil
}

// Get performs a GET request
func (c apiClient) Get(url string, rspData interface{}) (*http.Response, error) {
	r := Request{Method: http.MethodGet, URL: url}
	return c.Request(&r, rspData)
}

// Patch performs a PATCH request
func (c apiClient) Patch(url string, body map[string]interface{}, rspData interface{}) (*http.Response, error) {
	r := Request{Method: http.MethodPatch, URL: url, Body: body}
	return c.Request(&r, rspData)
}
