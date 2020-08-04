package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/taxi"
)

const (
	healthOk       = "healthy"
	healthDegraded = "degraded"
)

// HealthProvider returns a map representing the
// application's health. The keys should be the names of
// the checks.
type HealthProvider interface {
	Health(ctx context.Context) map[string]error
}

// CheckResult represents the result of a single health check
type CheckResult struct {
	Passed bool   `json:"passed"`
	Error  string `json:"error"`
}

// HealthResponse represents the JSON response that is
// returned by HealthHandler. It is returned under the
// `health` key in the response body.
type HealthResponse struct {
	Hostname string                  `json:"hostname"`
	Revision string                  `json:"revision"`
	Status   string                  `json:"status"`
	Checks   map[string]*CheckResult `json:"checks"`
}

// HealthHandler returns a JSON response describing the
// application's health. If all of the checks returned
// by the HealthProvider have nil errors, a 200 response
// is returned, otherwise a 500 response is returned.
//
// This handler should be used to provide a readiness
// probe. That is, a probes that determines whether the
// service is ready to receive traffic. It should not be
// used as a liveness probe because the checks performed
// by the HealthProvider may test connectivity to
// downstream services. It is usually undesirable for a
// service to be restarted when one of its dependencies is
// unavailable.
type HealthHandler struct {
	Hostname string
	Revision string
	Provider HealthProvider
}

// ServeHTTP returns a JSON response with health information
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rsp := &HealthResponse{
		Hostname: h.Hostname,
		Revision: h.Revision,
		Checks:   make(map[string]*CheckResult),
	}

	healthy := true

	results := h.Provider.Health(r.Context())
	for name, err := range results {
		rsp.Checks[name] = &CheckResult{Passed: err == nil}
		if err != nil {
			rsp.Checks[name].Error = err.Error()
			healthy = false
		}
	}

	var httpStatus int

	if healthy {
		rsp.Status = healthOk
		httpStatus = http.StatusOK
	} else {
		rsp.Status = healthDegraded
		httpStatus = http.StatusInternalServerError
	}

	payload := struct {
		Health *HealthResponse `json:"health"`
	}{rsp}

	b, err := json.Marshal(&payload)
	if err != nil {
		_ = taxi.WriteError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(httpStatus)
	_, err = w.Write(b)
}
