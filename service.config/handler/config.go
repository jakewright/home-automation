package handler

import (
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/request"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/service.config/domain"
	"github.com/jakewright/home-automation/service.config/service"
)

// ConfigHandler exports the handlers for the endpoints
type ConfigHandler struct {
	ConfigService *service.ConfigService
	Config        *domain.Config
}

type readRequest struct {
	ServiceName string `json:"service_name"`
}

// ReadConfig returns the config for the given service
func (h *ConfigHandler) ReadConfig(w http.ResponseWriter, r *http.Request) {
	body := &readRequest{}
	if err := request.Decode(r, body); err != nil {
		response.WriteJSON(w, err)
		return
	}

	config, err := h.Config.Get(body.ServiceName)
	if err != nil {
		response.WriteJSON(w, err)
		return
	}

	response.WriteJSON(w, config)
}

// ReloadConfig reads the YAML file from disk and loads changes into memory
func (h *ConfigHandler) ReloadConfig(w http.ResponseWriter, r *http.Request) {
	reloaded, err := h.ConfigService.Reload()
	if err != nil {
		response.WriteJSON(w, err)
		return
	}

	msg := "Config reloaded"

	if !reloaded {
		msg += " (no changes made)"
	}

	response.WriteJSON(w, msg)
}
