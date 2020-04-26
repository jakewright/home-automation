package router

import (
	"net/http"
	"runtime/debug"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/libraries/go/util"
)

// Revision is the service's revision and should be
// set at build time to the current commit hash.
var Revision string

func panicRecovery(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if v := recover(); v != nil {
			stack := debug.Stack()
			err := errors.Wrap(v, errors.ErrPanic, "recovered from panic", map[string]string{
				"stack": string(stack),
			})
			response.WriteJSON(w, err)

			if util.IsProd() {
				slog.Error(err)
			} else {
				// Panicking is useful in dev for the pretty-printed stack trace in terminal
				panic(err)
			}
		}
	}()

	next(w, r)
}

func revision(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if Revision != "" {
		w.Header().Set("X-Revision", Revision)
	}

	next(w, r)
}

// PingResponse is returned when a GET request to /ping is made
type PingResponse struct {
	Ping string `json:"ping,omitempty"`
}

func ping(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == http.MethodGet && r.URL.Path == "/ping" {
		response.WriteJSON(w, &PingResponse{
			Ping: "pong",
		})
		return
	}

	next(w, r)
}
