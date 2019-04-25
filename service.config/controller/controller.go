package controller

import (
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/service.config/domain"
	"github.com/jakewright/home-automation/service.config/service"

	"github.com/gorilla/mux"
)

// Controller exports the handlers for the endpoints
type Controller struct {
	ConfigService *service.ConfigService
	Config        *domain.Config
}

// ReadConfig returns the config for the given service
func (c *Controller) ReadConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName, ok := vars["serviceName"]
	if !ok {
		err := errors.InternalService("service name not provided", nil)
		response.WriteJSON(w, err)
		return
	}

	config, err := c.Config.Get(serviceName)
	if err != nil {
		response.WriteJSON(w, err)
		return
	}

	response.WriteJSON(w, config)
}

// ReloadConfig reads the YAML file from disk and loads changes into memory
func (c *Controller) ReloadConfig(w http.ResponseWriter, r *http.Request) {
	reloaded, err := c.ConfigService.Reload()
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
