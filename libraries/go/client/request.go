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

type RestRequester interface {
	Request(request *Request, resultData interface{}) (int, error)
	Get(url string, resultData interface{}) (int, error)
	Patch(url string, body map[string]interface{}, resultData interface{}) (int, error)
}

type APIClient struct {
	base *url.URL
}

func New(base string) (*APIClient, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		base: u,
	}, nil
}

type Request struct {
	Method string
	URL    string
	Body   map[string]interface{}
}

func (c APIClient) Request(request *Request, response interface{}) (status int, _ error) {
	// Convert the body to a byte array
	jsonBody, err := json.Marshal(request.Body)
	if err != nil {
		return 0, err
	}

	// If the URL is relative, prepend the base URL
	u, err := url.Parse(request.URL)
	resolvedUrl := c.base.ResolveReference(u)

	fmt.Println(resolvedUrl.String())

	client := &http.Client{}
	req, err := http.NewRequest(request.Method, resolvedUrl.String(), bytes.NewReader(jsonBody))
	if err != nil {
		return 0, err
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	// Validate the content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "json") {
		return 0, fmt.Errorf("content type %s not supported", contentType)
	}

	// The response will be enveloped in a data field so unmarshal
	// it into that format and then discard the outer layer.
	envelopedResponse := struct {
		Data interface{}
	}{
		Data: response,
	}

	// Read the body into a byte array
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Decode the byte array into a map
	err = json.Unmarshal(body, &envelopedResponse)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func (c APIClient) Get(url string, resultData interface{}) (status int, _ error) {
	r := Request{Method: http.MethodGet, URL: url}
	return c.Request(&r, resultData)
}

func (c APIClient) Patch(url string, body map[string]interface{}, resultData interface{}) (status int, _ error) {
	r := Request{Method: http.MethodPatch, URL: url, Body: body}
	return c.Request(&r, resultData)
}
