package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/service.dmx/ola"

	"github.com/jakewright/home-automation/service.dmx/repository"

	"github.com/jakewright/home-automation/libraries/go/request"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

type DMXHandler struct {
	Repository *repository.DMXRepository
}

type updateRequest struct {
	DeviceID   string `json:"device_id"` // URL param
	RGB        string `json:"rgb"`
	Strobe     int    `json:"strobe"`
	Brightness int    `json:"brightness"`
}

func (h *DMXHandler) Read(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := mux.Vars(r)["device_id"]
	if !ok {

	}

	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}

	fixture := h.Repository.Find(deviceID)
	if fixture == nil {
		response.WriteJSON(w, errors.NotFound("Device '%s' not found", deviceID))
	}

	request.Decode(r, &fixture)

	ola.SetDMX()
}

func (h *DMXHandler) Update(w http.ResponseWriter, r *http.Request) {
	var body updateRequest
	if err := request.Decode(r, &body); err != nil {
		slog.Error("Failed to decode body: %v", err)
		response.WriteJSON(w, err)
		return
	}
}
