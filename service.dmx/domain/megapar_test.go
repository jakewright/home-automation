package domain

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"

	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
)

func TestMegaParProfile_SetProperties(t *testing.T) {
	t.Parallel()

	type fields struct {
		power      bool
		color      color.RGBA
		strobe     byte
		brightness byte
	}

	tests := []struct {
		name    string
		state   map[string]interface{}
		wantErr bool
		before  fields
		after   fields
		changed bool
	}{
		{
			name:  "empty state",
			state: map[string]interface{}{},
			before: fields{
				power:      true,
				color:      color.RGBA{10, 10, 10, 0},
				strobe:     100,
				brightness: 50,
			},
			after: fields{
				power:      true,
				color:      color.RGBA{10, 10, 10, 0},
				strobe:     100,
				brightness: 50,
			},
		},
		{
			name: "no change",
			state: map[string]interface{}{
				"brightness": 50,
			},
			before: fields{
				brightness: 50,
			},
			after: fields{
				brightness: 50,
			},
		},
		{
			name: "change all fields",
			state: map[string]interface{}{
				"rgb":        "#FF0000",
				"strobe":     100,
				"brightness": 50,
				"power":      true,
			},
			before: fields{},
			after: fields{
				power:      true,
				color:      color.RGBA{255, 0, 0, 255},
				strobe:     100,
				brightness: 50,
			},
			changed: true,
		},
		{
			name: "brightness out of bounds",
			state: map[string]interface{}{
				"brightness": 300,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &MegaParProfile{
				abstractFixture: abstractFixture{
					DeviceHeader: &deviceregistrydef.DeviceHeader{},
				},
				power:      tt.before.power,
				color:      tt.before.color,
				strobe:     tt.before.strobe,
				brightness: tt.before.brightness,
			}

			changed, err := f.SetProperties(tt.state)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.changed, changed)

			require.Equal(t, f.power, tt.after.power)
			require.Equal(t, f.color, tt.after.color)
			require.Equal(t, f.strobe, tt.after.strobe)
			require.Equal(t, f.brightness, tt.after.brightness)
		})
	}
}
