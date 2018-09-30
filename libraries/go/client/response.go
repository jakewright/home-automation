package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type response struct {
	Value interface{} `json:"data"`
}

type errorResponse struct {
	Message string `json:"message"`
	Error   string `json:"errors"`
}

// Respond returns a response to the client
func Respond(w http.ResponseWriter, status int, data interface{}) {
	payload := response{Value: data}
	writeResponse(w, status, payload)
}

// RespondError returns the given error to the client with an appropriate status code
func RespondError(w http.ResponseWriter, et ErrorType, err error) {
	log.Println(err)

	payload := errorResponse{
		Message: et.Message,
		Error:   err.Error(),
	}

	writeResponse(w, et.Status, payload)
}

func writeResponse(w http.ResponseWriter, status int, payload interface{}) {
	// Set the content type
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	rsp, err := json.Marshal(&payload)
	if err != nil {
		// We could call RespondError here in a recursive fashion, but if that also fails then
		// we'll get stuck in an infinite loop.
		replacementPayload := errorResponse{
			Message: "Error converting payload to JSON",
			Error:   err.Error(),
		}

		rsp, err = json.Marshal(&replacementPayload)

		// If we can't even convert the error to JSON, send a minimal response back.
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, string(rsp))
}
