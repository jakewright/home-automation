package taxi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

// ResponseWriter constructs and writes responses to RPC requests
type ResponseWriter interface {
	Write(w http.ResponseWriter, v interface{}) error
}

type responseWriter struct{}

func (*responseWriter) Write(w http.ResponseWriter, v interface{}) error {
	return WriteResponse(w, v)
}

// WriteResponse marshals the body to JSON and wraps it in a
// "data" field in a JSON response. If the body is an error, the
// error's string is put in an "error" field in a JSON response.
func WriteResponse(w http.ResponseWriter, v interface{}) error {
	status := http.StatusOK
	payload := struct {
		Error string
		Data  interface{}
	}{}

	switch t := v.(type) {
	case *oops.Error:
		status = t.HTTPStatus()
		payload.Error = t.GetMessage()
	case error:
		status = http.StatusInternalServerError
		payload.Error = t.Error()
	default:
		payload.Data = v
	}

	rsp, err := json.Marshal(&payload)
	if err != nil {
		w.Header().Set("Content-Type", contentTypeText)
		w.WriteHeader(http.StatusInternalServerError)

		msg := fmt.Sprintf("Failed to marshal response payload: %v", err)
		if _, err := fmt.Fprintf(w, msg); err != nil {
			return err
		}

		return err
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	if _, err := fmt.Fprint(w, string(rsp)); err != nil {
		return err
	}

	return nil
}
