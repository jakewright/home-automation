package handler

import (
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/request"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

type updateRequest struct {
	DeviceID   string `json:"device_id"` // URL param
	RGB        string `json:"rgb"`
	Strobe     int    `json:"strobe"`
	Brightness int    `json:"brightness"`
}

func Update(w http.ResponseWriter, r *http.Request) {
	var body updateRequest
	if err := request.Decode(r, &body); err != nil {
		slog.Error("Failed to decode body: %v", err)
		response.WriteJSON(w, err)
		return
	}

}
