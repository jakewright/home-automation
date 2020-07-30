package domain

import (
	"image/color"

	"github.com/jakewright/home-automation/libraries/go/device"
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/ptr"
)

// MegaParProfile is a light by ADJ
type MegaParProfile struct {
	baseFixture

	power      bool
	color      device.RGB
	colorMacro byte
	strobe     byte
	program    byte
	brightness byte
}

var _ Fixture = (*MegaParProfile)(nil)

// length returns the number of DMX values that this fixture occupies
func (f *MegaParProfile) length() int { return 7 }

// hydrate sets the internal state based on the given DMX values
func (f *MegaParProfile) hydrate(values []byte) error {
	if len(values) != f.length() {
		return oops.InternalService(
			"expected % values to hydrate MegaParProfile but received %d",
			f.length(), len(values),
		)
	}

	f.color = device.RGB{
		RGBA: color.RGBA{
			R: values[0],
			G: values[1],
			B: values[2],
			A: 0xff,
		},
	}

	f.colorMacro = values[3]
	f.strobe = values[4]
	f.program = values[5]
	f.brightness = values[6]

	f.power = f.brightness > 0

	return nil
}

// dmxValues returns the DMX values for this fixture only
func (f *MegaParProfile) dmxValues() []byte {
	var b byte
	if f.power {
		b = f.brightness
	}

	return []byte{f.color.R, f.color.G, f.color.B, f.colorMacro, f.strobe, f.program, b}
}

// SetProperties sets any properties that exist in the state map
func (f *MegaParProfile) SetProperties(m map[string]interface{}) error {
	props := &MegaParProfileProperties{}
	if err := props.unmarshal(m); err != nil {
		return err
	}

	if props.Power != nil {
		f.power = *props.Power
	}
	if props.Rgb != nil {
		f.color = *props.Rgb
	}
	if props.Strobe != nil {
		f.strobe = byte(*props.Strobe)
	}
	if props.Brightness != nil {
		f.brightness = byte(*props.Brightness)
		f.power = *props.Brightness > 0
	}

	return nil
}

// ToDevice returns a standard Device type for a MegaParProfile
func (f *MegaParProfile) ToDevice() *devicedef.Device {
	state := &MegaParProfileProperties{
		Brightness: ptr.Int64(int64(f.brightness)),
		Power:      &f.power,
		Rgb:        &f.color,
		Strobe:     ptr.Int64(int64(f.strobe)),
	}

	return &devicedef.Device{
		Id:             f.ID(),
		Name:           f.Name,
		Type:           f.Type,
		Kind:           f.Kind,
		ControllerName: f.ControllerName,
		Attributes:     f.Attributes,
		StateProviders: nil,
		State:          state.describe(),
	}
}

// Copy returns a copy of the fixture but with zero values for the state
func (f *MegaParProfile) Copy() Fixture {
	return &MegaParProfile{
		baseFixture: f.baseFixture,
	}
}
