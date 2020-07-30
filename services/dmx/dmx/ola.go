package dmx

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jakewright/patch"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/services/dmx/domain"
)

// OLAClient makes requests to an OLA server to get and set DMX values
// https://wiki.openlighting.org/index.php/OLA_JSON_API
type OLAClient struct {
	un     domain.UniverseNumber
	client *patch.Client
}

// Compile-time assertion that OLAClient implements the GetSetter interface
var _ GetSetter = (*OLAClient)(nil)

// NewOLAClient returns a new OLA client for the given host:port
func NewOLAClient(host string, port int, un domain.UniverseNumber) (*OLAClient, error) {
	u, err := url.Parse(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, oops.WithMessage(err, "failed to parse OLA URL")
	}

	return &OLAClient{
		un: un,
		client: patch.New(
			patch.WithTimeout(time.Second*10),
			patch.WithStatusValidator(func(status int) bool {
				return status == 200
			}),
			patch.WithBaseURL(u.String()),
			patch.WithEncoder(&patch.EncoderFormURL{}),
		),
	}, nil
}

// GetValues returns the current DMX values
func (o *OLAClient) GetValues(ctx context.Context) ([512]byte, error) {
	u := fmt.Sprintf("/get_dmx?u=%d", o.un)

	rspBody := &struct {
		DMX   []byte `json:"dmx"`
		Error string `json:"error"`
	}{}

	if rsp, err := o.client.Get(ctx, u, nil); err != nil {
		return [512]byte{}, oops.WithMessage(err, "failed to make request to %s", u)
	} else if err := rsp.DecodeJSON(rspBody); err != nil {
		return [512]byte{}, oops.WithMessage(err, "failed to decode OLA response")
	}

	if rspBody.Error != "" {
		return [512]byte{}, oops.InternalService("received error from OLA: %s", rspBody.Error)
	}

	if len(rspBody.DMX) > 512 {
		return [512]byte{}, oops.InternalService("too many DMX values returned (%d)", len(rspBody.DMX))
	}

	values := [512]byte{}
	copy(values[:], rspBody.DMX)

	return values, nil
}

// SetValues sets the universe's DMX values
func (o *OLAClient) SetValues(ctx context.Context, values [512]byte) error {
	// Values must be sent as a comma-separated list
	var dst []byte
	for i, v := range values {
		if i > 0 {
			dst = append(dst, ","...)
		}
		dst = strconv.AppendUint(dst, uint64(v), 10)
	}

	body := struct {
		Universe domain.UniverseNumber `form:"u"`
		Values   string                `form:"d"`
	}{
		Universe: o.un,
		Values:   string(dst),
	}

	if _, err := o.client.Post(ctx, "/set_dmx", body, nil); err != nil {
		return oops.WithMessage(err, "failed to make request to /set_dmx")
	}

	return nil
}
