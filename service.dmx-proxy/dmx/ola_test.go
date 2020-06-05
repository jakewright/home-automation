package dmx

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_args(t *testing.T) {
	t.Parallel()

	tests := []struct {
		universe int
		values   [512]byte
		want     []string
	}{
		{
			universe: 0,
			values:   [512]byte{255, 0, 0, 0, 0, 0, 255},
			want:     []string{"--universe", "0", "--dmx", "255,0,0,0,0,0,255"},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			got := args(tt.universe, tt.values)
			require.Equal(t, tt.want, got)
		})
	}
}
