package util

import (
	"fmt"
	"image/color"

	"github.com/jakewright/home-automation/libraries/go/errors"
)

// ParseHexColor turns a hexadecimal color code (e.g. #FBEE13)
// into a color.RGBA.
func ParseHexColor(s string) (c color.RGBA, err error) {
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
		err = errors.BadRequest("invalid length; must be 7 or 4 characters")
	}

	if err != nil {
		err = errors.WrapWithCode(err, errors.ErrBadRequest, "failed to parse hex color code", map[string]string{
			"color_code": s,
		})
	}

	return
}
