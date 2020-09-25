package util

import (
	"encoding/json"
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

func TestRGB_MarshalJSON(t *testing.T) {
	s := &struct {
		Purple *RGB `json:"purple"`
		Black  *RGB `json:"black"`
		White  *RGB `json:"white"`
		Empty  *RGB `json:"empty"`
	}{
		Purple: &RGB{R: 54, G: 14, B: 142},
		Black:  &RGB{},
		White:  &RGB{R: 255, G: 255, B: 255},
		Empty:  nil,
	}

	want := `
{
	"purple": "#360E8E",
	"black": "#000000",
	"white": "#FFFFFF",
	"empty": null
}`

	b, err := json.Marshal(s)
	require.NoError(t, err)
	require.JSONEq(t, want, string(b))
}

func TestRGB_UnmarshalJSON(t *testing.T) {
	j := `
{
	"purple": "#360E8E",
	"black": "#000000",
	"white": "#FFFFFF",
	"empty": null
}`

	s := &struct {
		Purple *RGB `json:"purple"`
		Black  *RGB `json:"black"`
		White  *RGB `json:"white"`
		Empty  *RGB `json:"empty"`
	}{}

	err := json.Unmarshal([]byte(j), s)
	require.NoError(t, err)
	require.Equal(t, &RGB{R: 54, G: 14, B: 142, A: 0xFF}, s.Purple)
	require.Equal(t, &RGB{A: 0xFF}, s.Black)
	require.Equal(t, &RGB{R: 255, G: 255, B: 255, A: 0xFF}, s.White)
	require.Nil(t, s.Empty)
}
