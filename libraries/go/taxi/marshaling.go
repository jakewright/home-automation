package taxi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

const (
	contentTypeJSON = "application/json; charset=UTF-8"
	contentTypeText = "text/plain"
)

// DecodeRequest unmarshals URL parameters and the JSON body
// of the given request into the value pointed to by v.
// It is exported because it might be useful, e.g. in middleware.
func DecodeRequest(r *http.Request, v interface{}) error {
	// This does a load of reflection to unmarshal a map into the type of v
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true,
		Result:           v,

		// Override the TagName to match the one used by the
		// encoding/json package so users of this function only
		// have to define a single tag on struct fields
		TagName: "json",
	})
	if err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to create decoder")
	}

	// Query parameters come out as a map[string][]string so we loop through them all
	// to remove the unnecessary slice if the parameter just has a single value
	paramSlices := r.URL.Query()
	params := map[string]interface{}{}
	for key, value := range paramSlices {
		switch len(value) {
		case 0:
			params[key] = nil
		case 1:
			params[key] = value[0]
		default:
			params[key] = value
		}
	}

	// Unmarshal query parameters
	if err := decoder.Decode(params); err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to decode query parameters")
	}

	// If there's no body, return early
	if r.Body == nil {
		return nil
	}

	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to read request body")
	}

	if len(body) == 0 {
		return nil
	}

	// Assume the body is JSON and unmarshal into v
	if err := json.Unmarshal(body, v); err != nil {
		return oops.Wrap(err, oops.ErrBadRequest, "failed to unmarshal request body")
	}

	return nil
}

// WriteSuccess writes the data to the ResponseWriter.
// A status code of 200 is set.
func WriteSuccess(w http.ResponseWriter, v interface{}) error {
	payload := struct {
		Data interface{} `json:"data"`
	}{
		Data: v,
	}

	rsp, err := json.Marshal(&payload)
	if err != nil {
		// Best effort attempt to respond
		err = oops.Wrap(err, oops.ErrInternalService, "failed to marshal payload to JSON")
		_ = WriteError(w, err)
		return err
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(rsp)

	return err
}

// WriteError writes the error to the ResponseWriter. If the
// error is an oops.Error, the status code is taken from that,
// otherwise a status code of 500 is set.
func WriteError(w http.ResponseWriter, err error) error {
	status := http.StatusInternalServerError
	payload := struct {
		Error string `json:"error"`
		Stack string `json:"stack,omitempty"`
	}{Error: err.Error()}

	if oerr, ok := err.(*oops.Error); ok {
		status = oerr.HTTPStatus()
		payload.Error = oerr.GetMessage()

		// See the comment on the Format() function of
		// github.com/pkg/errors.StackTrace for formatting options
		payload.Stack = fmt.Sprintf("%+v", oerr.StackTrace())
	}

	rsp, err := json.Marshal(&payload)
	if err != nil {
		// Best effort attempt to respond
		w.Header().Set("Content-Type", contentTypeText)
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, "Failed to marshal error response payload to JSON: %s", err)
		return oops.WithMessage(err, "failed to marshal response payload to JSON")
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	_, err = w.Write(rsp)

	return err
}

// decodeResponse unmarshals the response
// to an RPC into the value point to by v
func decodeResponse(rsp *http.Response, v interface{}) error {
	var wrapper struct {
		Error string          `json:"error"`
		Data  json.RawMessage `json:"data"`
	}

	defer func() { _ = rsp.Body.Close() }()
	if err := json.NewDecoder(rsp.Body).Decode(&wrapper); err != nil {
		return oops.WithMessage(
			err,
			"received %d %s from server but could not decode response",
			rsp.StatusCode,
			rsp.Status,
		)
	}

	if wrapper.Error != "" {
		return oops.FromHTTPStatus(rsp.StatusCode, wrapper.Error)
	}

	if v == nil {
		return nil
	}

	if err := json.Unmarshal(wrapper.Data, v); err != nil {
		return oops.WithMessage(err, "failed to decode data")
	}

	return nil
}
