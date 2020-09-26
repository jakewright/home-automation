package domain

import (
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
)

// MockFixture can be used in tests
type MockFixture struct {
	IDValue string
	UN      UniverseNumber
	Ofs     int
	Len     int
}

var _ Fixture = (*MockFixture)(nil)

// ID returns the ID
func (f *MockFixture) ID() string {
	return f.IDValue
}

// UniverseNumber returns the universe number
func (f *MockFixture) UniverseNumber() UniverseNumber {
	return f.UN
}

// SetProperties sets the properties
func (f *MockFixture) SetProperties(m map[string]interface{}) error {
	panic("implement me")
}

// offset returns the offset
func (f *MockFixture) offset() int {
	return f.Ofs
}

// length returns the length
func (f *MockFixture) length() int {
	return f.Len
}

// hydrate is not implemented
func (f *MockFixture) hydrate(values []byte) error {
	panic("implement me")
}

// dmxValues is not implemented
func (f *MockFixture) dmxValues() []byte {
	panic("implement me")
}

// setHeader is not implemented
func (f *MockFixture) setHeader(header *devicedef.Header) error {
	panic("implement me")
}
