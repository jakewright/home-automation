package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"

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
func (f *MegaParProfile) SetProperties(data []byte) (bool, error) {
	var properties struct {
		RGB        string `json:"rgb"`
		Strobe     *byte  `json:"strobe"`
		Brightness *byte  `json:"brightness"`
		Power      *bool  `json:"power"`
	}

	if err := json.Unmarshal(data, &properties); err != nil {
		return false, err
	}

	var c color.RGBA
	if properties.RGB != "" {
		var err error
		if c, err = util.ParseHexColor(properties.RGB); err != nil {
			return false, err
		}
	}

	oldState := f.DMXValues()

	if properties.Power != nil {
		f.power = *properties.Power
	}
	if properties.Brightness != nil {
		f.brightness = *properties.Brightness
	}
	if properties.RGB != "" {
		f.color = c
	}
	if properties.Strobe != nil {
		f.strobe = *properties.Strobe
	}

	equal := bytes.Equal(oldState, f.DMXValues())
	return !equal, nil
}

// MarshalJSON returns the standard home-automation JSON encoding of the device
func (f *MegaParProfile) MarshalJSON() ([]byte, error) {
	rgb := fmt.Sprintf("#%02X%02X%02X", f.color.R, f.color.G, f.color.G)

	return json.Marshal(&struct {
		*DeviceHeader
		State map[string]interface{} `json:"state"`
	}{
		DeviceHeader: f.DeviceHeader,
		State: map[string]interface{}{
			"power": map[string]interface{}{
				"type":  "bool",
				"value": f.power,
			},
			"brightness": map[string]interface{}{
				"type":          "int",
				"min":           0,
				"max":           255,
				"interpolation": "continuous",
				"value":         f.brightness,
			},
			"rgb": map[string]interface{}{
				"type":  "rgb",
				"value": rgb,
			},
			"strobe": map[string]interface{}{
				"type":          "int",
				"min":           0,
				"max":           255,
				"interpolation": "continuous",
				"value":         f.strobe,
			},
		},
	})
}
