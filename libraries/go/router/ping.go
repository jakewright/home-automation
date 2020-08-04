package router

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/taxi"
)

// PingResponse is returned by PingHandler
type PingResponse struct {
	Ping string `json:"ping,omitempty"`
}

// PingHandler returns a simple pong response
func PingHandler(_ context.Context, _ taxi.Decoder) (interface{}, error) {
	return &PingResponse{
		Ping: "pong",
	}, nil
}
