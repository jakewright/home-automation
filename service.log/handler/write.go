package handler

import (
	"net/http"
	"time"

	"github.com/jakewright/home-automation/libraries/go/network"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

type writeRequest struct {
	Timestamp time.Time
	Severity  slog.Severity
	Message   string
	Metadata  map[string]string
}

// HandleWrite writes a slog line for testing purposes
func (h *Handler) HandleWrite(w http.ResponseWriter, r *http.Request) {
	body := writeRequest{}
	if err := network.DecodeRequest(r, &body); err != nil {
		network.WriteJSONResponse(w, err)
		return
	}

	if slog.DefaultLogger == nil {
		network.WriteJSONResponse(w, oops.InternalService("Default logger is nil"))
		return
	}

	if body.Timestamp.IsZero() {
		body.Timestamp = time.Now()
	}

	if int(body.Severity) == 0 {
		body.Severity = slog.InfoSeverity
	}

	if body.Message == "" {
		body.Message = "This is a log event"
	}

	if len(body.Metadata) == 0 {
		body.Metadata = map[string]string{"foo": "bar"}
	}

	event := &slog.Event{
		Timestamp: body.Timestamp,
		Severity:  body.Severity,
		Message:   body.Message,
		Metadata:  body.Metadata,
	}

	slog.DefaultLogger.Log(event)

	network.WriteJSONResponse(w, event)
}
