package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/libraries/go/taxi"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
	dmxdef "github.com/jakewright/home-automation/service.dmx/def"
	"github.com/jakewright/home-automation/service.dmx/dmx"
	"github.com/jakewright/home-automation/service.dmx/domain"
	"github.com/jakewright/home-automation/service.dmx/repository"
)

func TestDMXHandler_Update(t *testing.T) {
	// Setup the Firehose mock
	firehose.DefaultClient = &firehose.MockClient{}

	// Create a fixture
	f, err := domain.NewFixture(&deviceregistrydef.DeviceHeader{
		Id:             "fixture 1",
		Name:           "Fixture 1",
		Type:           "dmx",
		Kind:           "dmx",
		ControllerName: "service.dmx",
		Attributes: map[string]interface{}{
			"fixture_type": "mega_par_profile",
			"universe":     float64(1),
			"offset":       float64(0),
		},
	})
	require.NoError(t, err)

	// Set the fixture's initial state
	err = f.SetProperties(map[string]interface{}{
		"power":      false,
		"rgb":        "#FF0000",
		"strobe":     float64(0),
		"brightness": float64(0),
	})
	require.NoError(t, err)

	// Create a repository with the fixture
	repo := repository.New(f)

	client := dmx.NewClient()
	getSetter := &dmx.MockGetSetter{}
	client.AddGetSetter(1, getSetter)

	// Create the controller
	c := &Controller{
		Repository: repo,
		Client:     client,
	}

	d := &taxi.MockClient{Handler: NewRouter(c)}
	dmx := dmxdef.NewClient(d)

	rsp, err := dmx.UpdateDevice(context.Background(), &dmxdef.UpdateDeviceRequest{
		DeviceId: "fixture 1",
		State: map[string]interface{}{
			"brightness": float64(100),
			"rgb":        "#00FF00",
			"strobe":     float64(50),
		},
	}).Wait()

	require.NoError(t, err)

	require.Equal(t, "fixture 1", rsp.Device.Id)
	require.Equal(t, true, rsp.Device.State["power"].Value)
	require.EqualValues(t, 100, rsp.Device.State["brightness"].Value)
	require.EqualValues(t, "#00FF00", rsp.Device.State["rgb"].Value)
	require.EqualValues(t, 50, rsp.Device.State["strobe"].Value)

	expectedValues := [512]byte{0, 255, 0, 0, 50, 0, 100}
	require.Equal(t, expectedValues, getSetter.Values)
}
