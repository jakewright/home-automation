package domain

import (
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/util"
	dmxdef "github.com/jakewright/home-automation/services/dmx/def"
)

// MegaParProfile is a light by ADJ
type MegaParProfile struct {
	baseFixture

	power      bool
	color      util.RGB
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

	f.color = util.RGB{
		R: values[0],
		G: values[1],
		B: values[2],
		A: 0xff,
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

// ApplyState sets any properties that exist in the state map
func (f *MegaParProfile) ApplyState(p *dmxdef.MegaParProfileState) {
	if p == nil {
		return
	}

	if power, ok := p.GetPower(); ok {
		f.power = power
	}

	if color, ok := p.GetColor(); ok {
		f.color = color
	}

	if strobe, ok := p.GetStrobe(); ok {
		f.strobe = strobe
	}

	if brightness, ok := p.GetBrightness(); ok {
		f.brightness = brightness
		f.power = brightness > 0
	}
}

// State returns the current state of the device's properties
func (f *MegaParProfile) State() *dmxdef.MegaParProfileState {
	return (&dmxdef.MegaParProfileState{}).
		SetPower(f.power).
		SetBrightness(f.brightness).
		SetColor(f.color).
		SetStrobe(f.strobe)
}

// Copy returns a copy of the fixture but with zero values for the state
func (f *MegaParProfile) Copy() Fixture {
	return &MegaParProfile{
		baseFixture: f.baseFixture,
	}
}
