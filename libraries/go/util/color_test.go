package util

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHexToColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		hex  string
		want color.RGBA
	}{
		{
			hex:  "#FF0000",
			want: color.RGBA{255, 0, 0, 255},
		},
		{
			hex:  "#FBEE13",
			want: color.RGBA{251, 238, 19, 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			got, err := HexToColor(tt.hex)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestColorToHex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		c    color.RGBA
		want string
	}{
		{
			c:    color.RGBA{255, 0, 0, 255},
			want: "#FF0000",
		},
		{
			c:    color.RGBA{251, 238, 19, 255},
			want: "#FBEE13",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := ColorToHex(tt.c)
			require.Equal(t, tt.want, got)
		})
	}
}
