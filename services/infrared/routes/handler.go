package routes

import (
	"context"

	infrareddef "github.com/jakewright/home-automation/services/infrared/def"
	"github.com/jakewright/home-automation/services/infrared/ir"
	"github.com/jakewright/home-automation/services/infrared/repository"
)

type executor interface {
	Execute(context.Context, []ir.Instruction) error
}

type Controller struct {
	Repository *repository.DeviceRepository
	IR         executor
}

func (c *Controller) HandleGetDevice(r *Request, body *infrareddef.GetDeviceRequest) *dmx
