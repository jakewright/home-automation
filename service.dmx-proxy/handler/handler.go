package handler

import "github.com/jakewright/home-automation/service.dmx-proxy/dmx"

// Handler handles requests
type Handler struct {
	Setter dmx.Setter
}
