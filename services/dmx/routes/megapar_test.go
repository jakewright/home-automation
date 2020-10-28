package routes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/libraries/go/util"
	dmxdef "github.com/jakewright/home-automation/services/dmx/def"
	"github.com/jakewright/home-automation/services/dmx/dmx"
	"github.com/jakewright/home-automation/services/dmx/domain"
	"github.com/jakewright/home-automation/services/dmx/repository"
)

func TestController_Update(t *testing.T) {
	// Create a fixture
	f, err := domain.NewFixture((&devicedef.Header{}).
		SetId("fixture 1").
		SetName("Fixture 1").
		SetType("dmx").
		SetKind("dmx").
		SetControllerName("service.dmx").
		SetAttributes(map[string]interface{}{
			"fixture_type": "mega_par_profile",
			"universe":     float64(1),
			"offset":       float64(0),
		}),
	)
	require.NoError(t, err)

	megaParProfile, ok := f.(*domain.MegaParProfile)
	require.True(t, ok)

	// Set the fixture's initial state
	megaParProfile.ApplyState((&dmxdef.MegaParProfileState{}).
		SetPower(false).
		SetColor(util.RGB{R: 255}).
		SetStrobe(0).
		SetBrightness(0),
	)

	// Create a repository with the fixture
	repo := repository.New(f)

	client := dmx.NewClient()
	getSetter := &dmx.MockGetSetter{}
	client.AddGetSetter(1, getSetter)

	// Create the controller
	c := &Controller{
		Repository: repo,
		Client:     client,
		Publisher:  &firehose.MockClient{},
	}

	rsp, err := c.UpdateMegaParProfile(context.Background(), (&dmxdef.UpdateMegaParProfileRequest{
		State: (&dmxdef.MegaParProfileState{}).
			SetBrightness(100).
			SetColor(util.RGB{G: 255}).
			SetStrobe(50),
	}).SetDeviceId("fixture 1"))

	require.NoError(t, err)

	require.Equal(t, "fixture 1", rsp.Header.GetId())

	power, set := rsp.State.GetPower()
	require.Equal(t, true, set)
	require.Equal(t, true, power)

	brightness, set := rsp.State.GetBrightness()
	require.Equal(t, true, set)
	require.Equal(t, 100, brightness)

	color, set := rsp.State.GetColor()
	require.Equal(t, true, set)
	require.Equal(t, util.RGB{G: 255}, color)

	strobe, set := rsp.State.GetStrobe()
	require.Equal(t, true, set)
	require.Equal(t, 50, strobe)

	expectedValues := [512]byte{0, 255, 0, 0, 50, 0, 100}
	require.Equal(t, expectedValues, getSetter.Values)
}
