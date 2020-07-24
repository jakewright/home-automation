package device

import (
	"context"
	"image/color"

	"github.com/jakewright/home-automation/libraries/go/util"
)

//go:generate jrpc device.def

// Useful constants
const (
	InterpolationDiscrete   = "discrete"
	InterpolationContinuous = "continuous"
	TypeBool                = "bool"
	TypeInt                 = "int"
	TypeString              = "string"
	TypeRGB                 = "rgb"
)

// RGB wraps color.RGBA but provides marshaling
// functions to marshal to/from hex color codes
type RGB struct {
	color.RGBA
}

// UnmarshalText reads a hex value (e.g. #FF0000)
func (R *RGB) UnmarshalText(text []byte) error {
	c, err := util.HexToColor(string(text))
	if err != nil {
		return err
	}
	R.RGBA = c
	return nil
}

// MarshalText returns a hex value (e.g. #FF0000)
func (R *RGB) MarshalText() (text []byte, err error) {
	return []byte(util.ColorToHex(R.RGBA)), nil
}

// LoadProvidedState returns the state from a set of providers
func LoadProvidedState(ctx context.Context, deviceID string, providers []string) (map[string]interface{}, error) {
	state := make(map[string]interface{})

	for _, provider := range providers {
		//url := fmt.Sprintf("%s/provide-state?device_id=%s", provider, deviceID)

		_ = provider

		rsp := &struct {
			State map[string]interface{} `json:"state"`
		}{}

		//if _, err := rpc.Get(ctx, url, rsp); err != nil {
		//	return nil, err
		//}

		// Merge the new map
		for k, v := range rsp.State {
			state[k] = v
		}
	}

	return state, nil
}
