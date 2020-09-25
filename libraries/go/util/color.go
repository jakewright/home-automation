package util

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

// RGB is an alias for color.RGBA but with JSON marshaling
// methods that convert to & from hexadecimal color codes.
type RGB color.RGBA

// ToHex returns the hexadecimal representation of the
// RGB value, e.g. #FBEE13. The alpha value is ignored.
func (r *RGB) ToHex() string {
	return fmt.Sprintf("#%02X%02X%02X", r.R, r.G, r.B)
}

// MarshalJSON converts the RGB color into a hex value (e.g. "#FBEE13")
func (r *RGB) MarshalJSON() ([]byte, error) {
	if r == nil {
		return json.Marshal(nil)
	}

	b := make([]byte, 0, 7)
	b = append(b, '"')
	b = append(b, r.ToHex()...)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON converts a hex value (e.g. "#FF0000") to an RGB struct
func (r *RGB) UnmarshalJSON(b []byte) error {
	if r == nil {
		return oops.InternalService("cannot unmarshal into nil receiver")
	}

	// Unmarshalling JSON null is a no-op
	if string(b) == "null" {
		return nil
	}

	if b[0] != '"' || b[len(b)-1] != '"' {
		return oops.BadRequest("must be a quoted string to unmarshal color")
	}

	// Ignore the first and last char because they're the quotes
	c, err := HexToColor(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}

	r.R, r.G, r.B, r.A = c.R, c.G, c.B, c.A
	return nil
}

// HexToColor turns a hexadecimal color code
// (e.g. #FBEE13) into a color.RGBA.
func HexToColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = oops.BadRequest("invalid length; must be 7 or 4 characters")
	}

	if err != nil {
		err = oops.Wrap(err, oops.ErrBadRequest, "failed to parse hex color code", map[string]string{
			"color_code": s,
		})
	}

	return
}

// ColorToHex turns a color.RGBA into a hexadecimal color
// code e.g. #FBEE13. The alpha value is ignored.
func ColorToHex(c color.RGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}
