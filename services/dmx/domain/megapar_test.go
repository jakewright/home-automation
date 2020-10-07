package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/util"
	dmxdef "github.com/jakewright/home-automation/services/dmx/def"
)

func TestMegaParProfile_ApplyState(t *testing.T) {
	t.Parallel()

	type fields struct {
		power      bool
		color      util.RGB
		strobe     byte
		brightness byte
	}

	tests := []struct {
		name   string
		state  *dmxdef.MegaParProfileState
		before fields
		after  fields
	}{
		{
			name:  "empty state",
			state: nil,
			before: fields{
				power:      true,
				color:      util.RGB{R: 10, G: 10, B: 10, A: 0xff},
				strobe:     100,
				brightness: 50,
			},
			after: fields{
				power:      true,
				color:      util.RGB{R: 10, G: 10, B: 10, A: 0xff},
				strobe:     100,
				brightness: 50,
			},
		},
		{
			name:  "no change",
			state: (&dmxdef.MegaParProfileState{}).SetBrightness(50),
			before: fields{
				brightness: 50,
			},
			after: fields{
				brightness: 50,
			},
		},
		{
			name: "change all fields",
			state: (&dmxdef.MegaParProfileState{}).
				SetColor(util.RGB{R: 0xFF, A: 0xff}).
				SetStrobe(100).
				SetBrightness(50).
				SetPower(true),
			before: fields{},
			after: fields{
				power:      true,
				color:      util.RGB{R: 255, A: 0xff},
				strobe:     100,
				brightness: 50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &MegaParProfile{
				baseFixture: baseFixture{
					Header: &devicedef.Header{},
				},
				power:      tt.before.power,
				color:      tt.before.color,
				strobe:     tt.before.strobe,
				brightness: tt.before.brightness,
			}

			f.ApplyState(tt.state)

			require.Equal(t, tt.after.power, f.power)
			require.Equal(t, tt.after.color, f.color)
			require.Equal(t, tt.after.strobe, f.strobe)
			require.Equal(t, tt.after.brightness, f.brightness)
		})
	}
}
