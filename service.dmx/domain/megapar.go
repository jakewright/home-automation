package domain

import (
	"bytes"
	"encoding/json"
	"image/color"

	"github.com/jakewright/home-automation/libraries/go/device"
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/util"
)

// MegaParProfile is a light by ADJ
type MegaParProfile struct {
	*DeviceHeader

	power      bool
	color      color.RGBA
	colorMacro byte
	strobe     byte
	program    byte
	brightness byte
}

// ID returns the device ID
func (f *MegaParProfile) ID() string {
	return f.DeviceHeader.Id
}

// Offset returns the fixture's offset into the channel space
func (f *MegaParProfile) Offset() int {
	return f.Attributes.Offset
}

// DMXValues returns the DMX values for this fixture only
func (f *MegaParProfile) DMXValues() []byte {
	var b byte
	if f.power {
		b = f.brightness
	}

	return []byte{f.color.R, f.color.G, f.color.B, f.colorMacro, f.strobe, f.program, b}
}

// SetProperties unmarshals the []byte as JSON and sets
// any properties that exist in the resulting object.
func (f *MegaParProfile) SetProperties(state map[string]interface{}) (bool, error) {
	var properties struct {
		Power      *bool  `json:"power"`
		RGB        string `json:"rgb"`
		Strobe     *byte  `json:"strobe"`
		Brightness *byte  `json:"brightness"`
	}

	// The state map is the result of unmarshaling the JSON request
	// so all of the numbers end up being float64s. The easiest way
	// to turn these into the *byte types we want is to marshal back
	// to JSON and then unmarshal again. This deals with cases like
	// the number being too big to fit into a byte (uint8).

	b, err := json.Marshal(state)
	if err != nil {
		return false, errors.WithMessage(err, "failed to marshal state into JSON")
	}

	if err := json.Unmarshal(b, &properties); err != nil {
		return false, errors.WithMessage(err, "failed to unmarshal JSON")
	}

	var c color.RGBA
	if properties.RGB != "" {
		var err error
		if c, err = util.HexToColor(properties.RGB); err != nil {
			return false, errors.WithMessage(err, "failed to parse hex value")
		}
	}

	// Don't return any errors past this point otherwise the
	// in-memory fixture will be in an inconsistent state.

	oldState := f.DMXValues()

	if properties.Power != nil {
		f.power = *properties.Power
	}
	if properties.RGB != "" {
		f.color = c
	}
	if properties.Strobe != nil {
		f.strobe = *properties.Strobe
	}
	if properties.Brightness != nil {
		f.brightness = *properties.Brightness
		f.power = *properties.Brightness > 0
	}

	equal := bytes.Equal(oldState, f.DMXValues())
	return !equal, nil
}

// ToDef returns a standard Device type for a MegaParProfile
func (f *MegaParProfile) ToDef() *devicedef.Device {
	attributes := map[string]interface{}{
		"fixture_type": f.Attributes.FixtureType,
		"offset":       f.Attributes.Offset,
	}

	return &devicedef.Device{
		Id:             f.ID(),
		Name:           f.Name,
		Type:           f.Type,
		Kind:           f.Kind,
		ControllerName: f.ControllerName,
		Attributes:     attributes,
		StateProviders: nil,
		State: map[string]*devicedef.Property{
			"power": {
				Value: f.power,
				Type:  device.PropertyTypeBool,
			},
			"brightness": {
				Value:         f.brightness,
				Type:          device.PropertyTypeInt,
				Min:           0,
				Max:           255,
				Interpolation: device.InterpolationContinuous,
			},
			"rgb": {
				Value: util.ColorToHex(f.color),
				Type:  "rgb",
			},
			"strobe": {
				Value:         f.strobe,
				Type:          device.PropertyTypeInt,
				Min:           0,
				Max:           255,
				Interpolation: device.InterpolationContinuous,
			},
		},
	}
}
