package device

import (
	"context"
	"fmt"
	"image/color"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/rpc"
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

// BoolProperty returns a boolean property
func BoolProperty(value bool) *devicedef.Property {
	return &devicedef.Property{
		Value: value,
		Type:  TypeBool,
	}
}

// IntProperty returns an integer property
func IntProperty(value int, min, max float64, interpolation string) *devicedef.Property {
	return &devicedef.Property{
		Value:         value,
		Type:          TypeInt,
		Min:           &min,
		Max:           &max,
		Interpolation: interpolation,
	}
}

// Uint8Property returns a uint8 property
func Uint8Property(value uint8, min, max float64, interpolation string) *devicedef.Property {
	return &devicedef.Property{
		Value:         value,
		Type:          TypeInt,
		Min:           &min,
		Max:           &max,
		Interpolation: interpolation,
	}
}

// RGBProperty returns an RGB property
func RGBProperty(value color.RGBA) *devicedef.Property {
	return &devicedef.Property{
		Value: util.ColorToHex(value),
		Type:  TypeRGB,
	}
}

// StringProperty returns a string property
func StringProperty(value string) *devicedef.Property {
	return &devicedef.Property{
		Value: value,
		Type:  TypeString,
	}
}

// IntArg returns an integer argument
func IntArg(min, max float64, required bool) *devicedef.Arg {
	return &devicedef.Arg{
		Type:     TypeInt,
		Min:      &min,
		Max:      &max,
		Required: required,
	}
}

// StringArgWithOptions returns a string argument with a pick list of options
func StringArgWithOptions(options []*devicedef.Option, required bool) *devicedef.Arg {
	return &devicedef.Arg{
		Type:     TypeString,
		Options:  options,
		Required: required,
	}
}

// ValidateState returns an error if the given state does not conform to the specification
func ValidateState(state map[string]interface{}, device *devicedef.Device) error {
	spec := device.State

	for property, value := range state {
		def, ok := spec[property]
		if !ok {
			return oops.BadRequest("%q is not a valid property", property)
		}

		if err := validate(def.Type, def.Min, def.Max, def.Options, value); err != nil {
			return oops.WithMessage(err, "failed to validate %q", property, map[string]string{
				"property": property,
				"value":    fmt.Sprintf("%v", value),
			})
		}
	}

	return nil
}

// ValidateCommand returns an error if the command and its arguments do not conform to the specification
func ValidateCommand(command string, args map[string]interface{}, spec map[string]*devicedef.Command) error {
	cmd, ok := spec[command]
	if !ok {
		return oops.BadRequest("%q is not a valid command", command)
	}

	for argName, argDef := range cmd.Args {
		val, ok := args[argName]
		if !ok {
			if argDef.Required {
				return oops.BadRequest("%q is a required argument", argName)
			}

			continue
		}

		if err := validate(argDef.Type, argDef.Min, argDef.Max, argDef.Options, val); err != nil {
			return oops.WithMetadata(err, map[string]string{
				"arg": argName,
				"val": fmt.Sprintf("%v", val),
			})
		}
	}

	return nil
}

func validate(t string, min, max *float64, options []*devicedef.Option, v interface{}) error {
	switch t {
	case TypeBool:
		if _, ok := v.(bool); !ok {
			return oops.BadRequest("expected type %s but got %T", TypeBool, v)
		}

		return nil

	case TypeInt:
		// Expect float64 because encoding/json
		// unmarshals JSON numbers into float 64
		f, ok := v.(float64)
		if !ok {
			return oops.BadRequest("expected type %s but got %T", TypeInt, v)
		}

		if min != nil {
			if f < *min {
				return oops.BadRequest("value of %f is lower than minimum %f", f, *min)
			}
		}

		if max != nil {
			if f > *max {
				return oops.BadRequest("value of %f is higher than maximum %f", f, *max)
			}
		}

	case TypeString:
		s, ok := v.(string)
		if !ok {
			return oops.BadRequest("expected type %s but got %T", TypeString, v)
		}

		if len(options) > 0 {
			valid := false
			for _, op := range options {
				if s == op.Value {
					valid = true
					break
				}
			}
			if !valid {
				return oops.BadRequest("invalid value %q", s)
			}
		}
	}

	// todo: support other types

	return nil
}

// LoadProvidedState returns the state from a set of providers
func LoadProvidedState(ctx context.Context, deviceID string, providers []string) (map[string]interface{}, error) {
	state := make(map[string]interface{})

	for _, provider := range providers {
		url := fmt.Sprintf("%s/provide-state?device_id=%s", provider, deviceID)

		rsp := &struct {
			State map[string]interface{} `json:"state"`
		}{}

		if _, err := rpc.Get(ctx, url, rsp); err != nil {
			return nil, err
		}

		// Merge the new map
		for k, v := range rsp.State {
			state[k] = v
		}
	}

	return state, nil
}
