package routes

import (
	"home-automation/libraries/go/router"

	"home-automation/service.log/dao"
)

type Controller struct {
	Repository *dao.LogDAO
}

func (c *Controller) RegisterRoutes() {
	router.Get("/read", c.handleReadLogs)
	router.Get("/ws", c.handleWebSocket)
}
