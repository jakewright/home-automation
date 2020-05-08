package domain

import (
	"github.com/jakewright/home-automation/libraries/go/errors"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

type abstractFixture struct {
	*deviceregistrydef.DeviceHeader
	offset int
}

// SetHeader sets the fixture's header and pulls the offset out of the attributes
func (f *abstractFixture) SetHeader(h *deviceregistrydef.DeviceHeader) error {
	offset, ok := h.Attributes["offset"].(float64)
	if !ok {
		return errors.PreconditionFailed("offset not found in %s device header", h.Id)
	}

	f.DeviceHeader = h
	f.offset = int(offset)
	return nil
}

// ID returns the device ID
func (f *abstractFixture) ID() string {
	return f.DeviceHeader.Id
}

// Offset returns the fixture's offset into the channel space
func (f *abstractFixture) Offset() int {
	return f.offset
}
