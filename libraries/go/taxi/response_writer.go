package taxi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

type responseWriter struct {
	logFunc func(format string, v ...interface{})
}

func (rw *responseWriter) log(format string, v ...interface{}) {
	if rw.logFunc == nil {
		return
	}

	rw.logFunc(format, v...)
}

func (rw *responseWriter) writeResponse(w http.ResponseWriter, v interface{}) {
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
		if _, err = fmt.Fprintf(w, msg); err != nil {
			rw.log("Failed to write response: %v", err)
		}

		rw.log(msg)
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	if _, err := fmt.Fprint(w, string(rsp)); err != nil {
		rw.log("Failed to write response: %v", err)
	}
}
