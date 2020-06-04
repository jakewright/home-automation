package handler

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/libraries/go/test"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
	dmxproxydef "github.com/jakewright/home-automation/service.dmx-proxy/def"
	dmxdef "github.com/jakewright/home-automation/service.dmx/def"
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
			"offset":       float64(7), // force float64 to replicate what json.Unmarshal
		},
	})
	require.NoError(t, err)

	// Set the fixture's initial state
	_, err = f.SetProperties(map[string]interface{}{
		"power":      false,
		"rgb":        "#FF0000",
		"strobe":     float64(0),
		"brightness": float64(0),
	})
	require.NoError(t, err)

	// Create a universe and add the fixture
	u := universe.New(1)
	u.AddFixture(f)

	m, ctx := test.NewMock(t)
	defer m.Stop()

	expectedDMXValues := [512]byte{0, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 50, 0, 100}
	m.ExpectOne(&dmxproxydef.SetRequest{
		Universe: 1,
		Values:   expectedDMXValues[:],
	}).RespondWith(&dmxproxydef.SetResponse{})

	// Create the controller
	h := &Handler{
		Universe: u,
	}

	rsp, err := h.UpdateDevice(&request{Context: ctx},
		&dmxdef.UpdateDeviceRequest{
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

	m.RunAssertions()
}
