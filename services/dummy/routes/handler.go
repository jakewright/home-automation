package routes

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/slog"
	def "github.com/jakewright/home-automation/services/dummy/def"
)

// Controller handles requests
type Controller struct{}

// Log emits a log line
func (c Controller) Log(ctx context.Context, body *def.LogRequest) (*def.LogResponse, error) {
	slog.Infof("This is a log line", map[string]string{
		"foo": "bar",
	})
	return &def.LogResponse{}, nil
}

// Panic panics (it will be recovered by the framework)
func (c Controller) Panic(ctx context.Context, body *def.PanicRequest) (*def.PanicResponse, error) {
	panic("Panic! at the Handler")
}
