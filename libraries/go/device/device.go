package device

import (
	"context"
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
