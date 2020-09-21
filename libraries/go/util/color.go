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

// MarshalJSON converts the RGB color into a hex value (e.g. #FBEE13)
func (r *RGB) MarshalJSON() ([]byte, error) {
	if r == nil {
		return json.Marshal(nil)
	}

	return []byte(ColorToHex(color.RGBA(*r))), nil
}

// UnmarshalJSON converts a hex value (e.g. #FF0000) to an RGB struct
func (r *RGB) UnmarshalJSON(b []byte) error {
	if r == nil {
		return oops.InternalService("cannot unmarshal into nil receiver")
	}

	c, err := HexToColor(string(b))
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
