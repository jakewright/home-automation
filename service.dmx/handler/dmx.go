package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/service.dmx/domain"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/service.dmx/ola"

	"github.com/jakewright/home-automation/libraries/go/response"
)

type DMXHandler struct {
	Universe *domain.Universe
}

func (h *DMXHandler) Read(w http.ResponseWriter, r *http.Request) {
	// Get the device ID from the route params
	deviceID, ok := mux.Vars(r)["device_id"]
	if !ok {
		response.WriteJSON(w, errors.BadRequest("device_id not set in route params"))
		return
	}

	fixture := h.Universe.Find(deviceID)
	if fixture == nil {
		response.WriteJSON(w, errors.NotFound("Device '%s' not found", deviceID))
		return
	}

	response.WriteJSON(w, fixture)
}

func (h *DMXHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Get the device ID from the route params
	deviceID, ok := mux.Vars(r)["device_id"]
	if !ok {
		response.WriteJSON(w, errors.BadRequest("device_id not set in route params"))
		return
	}

	fixture := h.Universe.Find(deviceID)
	if fixture == nil {
		response.WriteJSON(w, errors.NotFound("Device '%s' not found", deviceID))
		return
	}

	// Read the body of the request
	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteJSON(w, errors.Wrap(err, "failed to read request body"))
		return
	}

	changed, err := fixture.SetProperties(body)
	if err != nil {
		response.WriteJSON(w, errors.Wrap(err, "failed to update fixture"))
		return
	}

	if err := ola.SetDMX(h.Universe.Number, h.Universe.DMXValues()); err != nil {
		response.WriteJSON(w, errors.Wrap(err, "failed to set DMX values"))
		return
	}

	if changed {
		if err := firehose.Publish("device-state-changed."+deviceID, fixture); err != nil {
			response.WriteJSON(w, errors.Wrap(err, "failed to emit event"))
			return
		}
	}

	response.WriteJSON(w, fixture)
}
