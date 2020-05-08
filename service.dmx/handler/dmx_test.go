package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jakewright/home-automation/libraries/go/firehose"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
	dmxdef "github.com/jakewright/home-automation/service.dmx/def"
	"github.com/jakewright/home-automation/service.dmx/dmx"
	"github.com/jakewright/home-automation/service.dmx/domain"
	"github.com/jakewright/home-automation/service.dmx/universe"
)

func TestDMXHandler_Update(t *testing.T) {
	// Setup the Firehose mock
	firehose.DefaultClient = &firehose.MockClient{}

	// Create a fixture
	f, err := domain.NewFixtureFromDeviceHeader(&deviceregistrydef.DeviceHeader{
		Id: "fixture 1",
		Attributes: map[string]interface{}{
			"fixture_type": domain.FixtureTypeMegaParProfile,
			"offset":       7,
		},
	})
	require.NoError(t, err)

	// Set the fixture's initial state
	_, err = f.SetProperties(map[string]interface{}{
		"power":      false,
		"rgb":        "#FF0000",
		"strobe":     float64(0), // force float64 to replicate what json.Unmarshal
		"brightness": float64(0),
	})
	require.NoError(t, err)

	// Create a universe and add the fixture
	u := universe.New(1)
	u.AddFixture(f)

	// Create a mock DMX setter
	s := &dmx.Mock{}

	// Create the controller
	h := &DMXController{
		Universe: u,
		Setter:   s,
	}

	rsp, err := h.Update(&Request{
		Context: context.Background(),
	}, &dmxdef.UpdateDeviceRequest{
		DeviceId: "fixture 1",
		State: map[string]interface{}{
			"brightness": float64(100),
			"rgb":        "#00FF00",
			"strobe":     float64(50),
		},
	})
	require.NoError(t, err)

	require.Equal(t, "fixture 1", rsp.Device.Id)
	require.Equal(t, true, rsp.Device.State["power"].Value)
	require.Equal(t, byte(100), rsp.Device.State["brightness"].Value)
	require.Equal(t, "#00FF00", rsp.Device.State["rgb"].Value)
	require.Equal(t, byte(50), rsp.Device.State["strobe"].Value)

	expectedDMXValues := [512]byte{0, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 50, 0, 100}
	require.Equal(t, 1, s.Universe)
	require.Equal(t, expectedDMXValues, s.Values, "expected %v\ngot      %v", expectedDMXValues, s.Values)
}
