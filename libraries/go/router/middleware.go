package router

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/jakewright/home-automation/libraries/go/environment"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/libraries/go/taxi"
)

// Revision is the service's revision and should be
// set at build time to the current commit hash.
var Revision string

func panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				stack := debug.Stack()
				err := oops.Wrap(v, oops.ErrInternalService, "recovered from panic", map[string]string{
					"stack": string(stack),
				})
				if err := taxi.WriteResponse(w, err); err != nil {
					slog.Errorf("Failed to write response: %v", err)
				}

				if environment.IsProd() {
					slog.Error(err)
				} else {
					// Panicking is useful in dev for the pretty-printed stack trace in terminal
					panic(err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func revision(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Revision != "" {
			w.Header().Set("X-Revision", Revision)
		}

		next.ServeHTTP(w, r)
	})
}

// PingResponse is returned when a GET request to /ping is made
type PingResponse struct {
	Ping string `json:"ping,omitempty"`
}

func pingHandler(_ context.Context, _ taxi.Decoder) (interface{}, error) {
	return &PingResponse{
		Ping: "pong",
	}, nil
}
