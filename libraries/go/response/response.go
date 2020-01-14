package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/slog"

	"github.com/jakewright/home-automation/libraries/go/errors"
)

type response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Write returns a raw response to the client
func Write(w http.ResponseWriter, buf bytes.Buffer) {
	if _, err := buf.WriteTo(w); err != nil {
		log.Println("Failed to write response", err)
	}
}

// WriteJSON returns a response to the client. If data is an error,
// a status code is inferred. Otherwise, data is marshaled to JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) {
	status := http.StatusOK
	payload := response{}

	if e, ok := data.(*errors.Error); ok {
		status = e.HTTPStatus()
		payload.Message = e.Error()
	} else if e, ok := (data).(error); ok {
		status = http.StatusInternalServerError
		payload.Message = e.Error()
	} else {
		payload.Data = data
	}

	writeJSON(w, status, payload)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	rsp, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Errorf("Failed to marshal response payload: %v", err)
		return
	}

	w.WriteHeader(status)
	if _, err := fmt.Fprint(w, string(rsp)); err != nil {
		slog.Errorf("Failed to write response: %v", err)
	}
}
