package response

import (
	"encoding/json"
	"fmt"
	"home-automation/libraries/go/errors"
	"log"
	"net/http"
)

type response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Write(w http.ResponseWriter, body string) {
	if _, err := fmt.Fprint(w, body); err != nil {
		log.Println("Failed to write response", err)
	}
}

// WriteJSON returns a response to the client
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
		log.Println(err)
		return
	}

	w.WriteHeader(status)
	if _, err := fmt.Fprint(w, string(rsp)); err != nil {
		log.Println("Failed to write response", err)
	}
}
