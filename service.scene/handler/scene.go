package handler

import (
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/request"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.scene/domain"
)

func HandleCreateScene(w http.ResponseWriter, r *http.Request) {
	scene := &domain.Scene{}
	if err := request.Decode(r, &scene); err != nil {
		slog.Error("Failed to decode body: %v", err)
		response.WriteJSON(w, err)
		return
	}

	if err := database.Create(scene).Error; err != nil {

	}
}
