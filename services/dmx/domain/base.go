package domain

import (
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/oops"
)

type baseFixture struct {
	*devicedef.Header
	universeNumber UniverseNumber
	offsetValue    int
}

// setHeader sets the fixture's header and pulls the offset out of the attributes
func (f *baseFixture) setHeader(h *devicedef.Header) error {
	universeNumber, ok := h.Attributes["universe"].(float64)
	if !ok {
		return oops.PreconditionFailed("universe number not found in %s device header", h.Id)
	}

	offset, ok := h.Attributes["offset"].(float64)
	if !ok {
		return oops.PreconditionFailed("offset not found in %s device header", h.Id)
	}

	f.Header = h
	f.universeNumber = UniverseNumber(universeNumber)
	f.offsetValue = int(offset)
	return nil
}

// ID returns the device ID
func (f *baseFixture) ID() string { return f.Header.GetId() }

// UniverseNumber returns the fixture's universe number
func (f *baseFixture) UniverseNumber() UniverseNumber { return f.universeNumber }

// offset returns the fixture's offset into the channel space
func (f *baseFixture) offset() int { return f.offsetValue }
