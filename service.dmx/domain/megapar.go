package domain

import (
	"bytes"
	"encoding/json"
	"image/color"

	"github.com/jakewright/home-automation/libraries/go/device"
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/util"
)

// MegaParProfile is a light by ADJ
type MegaParProfile struct {
	abstractFixture

	power      bool
	color      color.RGBA
	colorMacro byte
	strobe     byte
	program    byte
	brightness byte
}

// DMXValues returns the DMX values for this fixture only
func (f *MegaParProfile) DMXValues() []byte {
	var b byte
	if f.power {
		b = f.brightness
	}

	return []byte{f.color.R, f.color.G, f.color.B, f.colorMacro, f.strobe, f.program, b}
}

// SetProperties sets any properties that exist in the state map
func (f *MegaParProfile) SetProperties(state map[string]interface{}) (bool, error) {
	if err := device.ValidateState(state, f.ToDef()); err != nil {
		return false, err
	}

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
		return false, oops.WithMessage(err, "failed to marshal state into JSON")
	}

	if err := json.Unmarshal(b, &properties); err != nil {
		return false, oops.WithMessage(err, "failed to unmarshal JSON")
	}

	var c color.RGBA
	if properties.RGB != "" {
		var err error
		if c, err = util.HexToColor(properties.RGB); err != nil {
			return false, oops.WithMessage(err, "failed to parse hex value")
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
	return &devicedef.Device{
		Id:             f.ID(),
		Name:           f.Name,
		Type:           f.Type,
		Kind:           f.Kind,
		ControllerName: f.ControllerName,
		Attributes:     f.Attributes,
		StateProviders: nil,
		State: map[string]*devicedef.Property{
			"power":      device.BoolProperty(f.power),
			"brightness": device.Uint8Property(f.brightness, 0, 255, device.InterpolationContinuous),
			"rgb":        device.RGBProperty(f.color),
			"strobe":     device.Uint8Property(f.strobe, 0, 255, device.InterpolationContinuous),
		},
	}
}

// Copy returns a copy of the fixture
func (f *MegaParProfile) Copy() Fixture {
	return &MegaParProfile{
		abstractFixture: f.abstractFixture, // note this is not a deep copy
		power:           f.power,
		color:           f.color,
		colorMacro:      f.colorMacro,
		strobe:          f.strobe,
		program:         f.program,
		brightness:      f.brightness,
	}
}
