package handler

import (
	"strconv"
	"strings"

	"github.com/jakewright/home-automation/libraries/go/oops"
	dmxproxydef "github.com/jakewright/home-automation/service.dmx-proxy/def"
	"github.com/jakewright/home-automation/service.dmx-proxy/dmx"
)

// Controller handles requests
type Controller struct {
	Setter dmx.Setter
}

// Set returns a handler that sets DMX values
func (c *Controller) Set(_ *request, body *dmxproxydef.SetRequest) (*dmxproxydef.SetResponse, error) {
	var values [512]byte
	copy(values[:], body.Values)

	valuesStr := make([]string, len(body.Values))
	for i, v := range body.Values {
		valuesStr[i] = strconv.Itoa(int(v))
	}

	if err := c.Setter.Set(int(body.Universe), values); err != nil {
		return nil, oops.WithMessage(err, "failed to set DMX values", map[string]string{
			"universe": strconv.Itoa(int(body.Universe)),
			"values":   strings.Join(valuesStr, ","),
		})
	}

	return &dmxproxydef.SetResponse{}, nil
}
