package routes

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"sync"

	"github.com/jakewright/home-automation/libraries/go/oops"
	def "github.com/jakewright/home-automation/services/lirc-proxy/def"
)

const (
	irSend   = "irsend"    // the LIRC executable
	sendOnce = "SEND_ONCE" // the LIRC command
)

// Controller handles requests
type Controller struct {
	mu sync.Mutex
}

// SendOnce sends the device and key to LIRC
func (c *Controller) SendOnce(
	ctx context.Context,
	body *def.SendOnceRequest,
) (*def.SendOnceResponse, error) {
	// Avoid concurrent handlers spawning multiple irsend
	// processes simultaneously. I'm not sure what would
	// happen in this situation but it's easy to avoid,
	// assuming there is only one instance of this service.
	c.mu.Lock()
	defer c.mu.Unlock()

	switch {
	case body.GetDevice() == "":
		return nil, oops.BadRequest(
			"field 'device' must not be empty",
		)
	case body.GetKey() == "":
		return nil, oops.BadRequest(
			"field 'key' must not be empty",
		)
	}

	args := []string{sendOnce, body.GetDevice(), body.GetKey()}
	cmd := exec.CommandContext(ctx, irSend, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	if err := cmd.Run(); err != nil {
		return nil, oops.WithMessage(
			err,
			"failed to run `%s %s`",
			irSend,
			strings.Join(args, " "),
		)
	}

	// Assume anything written to stderr is a bad thing
	if stderr.Len() > 0 {
		return nil, oops.InternalService(
			"%s wrote to stderr: %s",
			irSend,
			stderr.String(),
		)
	}

	return &def.SendOnceResponse{}, nil
}
