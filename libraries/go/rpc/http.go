package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

const defaultTimeout = 10 * time.Second

// HTTPClient is a high-level HTTP client for making requests to other services
type HTTPClient struct {
	// Base is the base URL used for relative requests
	Base string

	// ValidateStatus is a function that validates the status code of the HTTP response
	ValidateStatus func(int) bool

	// Envelope is the name of the data field in the response
	Envelope string

	httpClient *http.Client
}

// NewHTTPClient returns a new RPC HTTPClient
func NewHTTPClient(envelope string) (Client, error) {
	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	return &HTTPClient{
		ValidateStatus: func(status int) bool {
			return status >= 200 && status < 300
		},
		Envelope: envelope,

		httpClient: httpClient,
	}, nil
}

// Do performs the HTTP request. Relative URLs will have the base URL prepended. An error
// will be thrown if the response does not have a JSON content-type or if the status code is
// not valid. The entire response will be returned but the JSON will be unmarshalled into the second argument.
func (c *HTTPClient) Do(ctx context.Context, request *Request, v interface{}) (*Response, error) {
	if !validMethod(request.Method) {
		return nil, fmt.Errorf("method %q is not a valid HTTP method", request.Method)
	}

	// Convert the body to a byte array
	var jsonBody []byte
	var err error
	if request.Body != nil {
		jsonBody, err = json.Marshal(request.Body)
		if err != nil {
			return nil, err
		}
	}

	// If the URL is relative, prepend the base URL
	absURL := request.URL
	u, err := url.Parse(request.URL)
	if err != nil {
		return nil, err
	}
	if c.Base != "" && !u.IsAbs() {
		absURL = fmt.Sprintf("%s/%s", c.Base, strings.TrimLeft(request.URL, "/"))
	}

	// Construct the request
	req, err := http.NewRequest(request.Method, absURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	req = req.WithContext(ctx)

	// Make the request
	rawRsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// From this point on, all return values should return rsp, even if there's an error
	// so that the caller can see all of the information about the response.
	rsp := &Response{Response: rawRsp}

	// Read the body into a byte array
	defer func() { _ = rawRsp.Body.Close() }()
	body, err := ioutil.ReadAll(rawRsp.Body)
	if err != nil {
		return rsp, err
	}
	rsp.Body = body

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
			return rsp, oops.WithMessage(err, "failed to unmarshal json %s", innerBytes)
		}
	} else {
		if err = json.Unmarshal(body, &v); err != nil {
			return rsp, err
		}
	}

	return rsp, nil
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
