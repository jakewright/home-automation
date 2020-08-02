package router

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/jakewright/home-automation/libraries/go/bootstrap"
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
				if err := taxi.WriteError(w, err); err != nil {
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

const (
	healthOk       = "healthy"
	healthDegraded = "degraded"
)

type healthHandler struct {
	svc *bootstrap.Service
}

// ServeHTTP returns a JSON response with health information
func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type checkResult struct {
		Passed bool   `json:"passed"`
		Error  string `json:"error"`
	}

	type health struct {
		Status string                  `json:"status"`
		Checks map[string]*checkResult `json:"checks"`
	}

	payload := struct {
		Health *health `json:"health"`
	}{
		Health: &health{
			Checks: make(map[string]*checkResult),
		},
	}

	healthy := true
	results := h.svc.Healthy(r.Context())
	for name, err := range results {
		payload.Health.Checks[name] = &checkResult{Passed: err == nil}
		if err != nil {
			payload.Health.Checks[name].Error = err.Error()
			healthy = false
		}
	}

	var httpStatus int

	if healthy {
		payload.Health.Status = healthOk
		httpStatus = http.StatusOK
	} else {
		payload.Health.Status = healthDegraded
		httpStatus = http.StatusInternalServerError
	}

	rsp, err := json.Marshal(&payload)
	if err != nil {
		_ = taxi.WriteError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(httpStatus)
	_, err = w.Write(rsp)
}
