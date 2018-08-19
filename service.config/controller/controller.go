package controller

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"home-automation/libraries/go/client"
	"home-automation/service.config/domain"
	"home-automation/service.config/service"
	"net/http"
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
		err := errors.New("service name not provided")
		client.RespondError(w, client.ErrInternalService, err)
		return
	}

	config, err := c.Config.Get(serviceName)
	if err != nil {
		client.RespondError(w, client.ErrInternalService, err)
		return
	}

	client.Respond(w, fmt.Sprintf("Config for %q", serviceName), config, http.StatusOK)
}

// ReloadConfig reads the YAML file from disk and loads changes into memory
func (c *Controller) ReloadConfig(w http.ResponseWriter, r *http.Request) {
	reloaded, err := c.ConfigService.Reload()
	if err != nil {
		client.RespondError(w, client.ErrInternalService, err)
		return
	}

	msg := "Config reloaded"

	if !reloaded {
		msg += " (no changes made)"
	}

	client.Respond(w, msg, nil, http.StatusOK)
}
