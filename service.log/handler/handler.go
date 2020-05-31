package handler

import (
	"github.com/jakewright/home-automation/service.log/repository"
	"github.com/jakewright/home-automation/service.log/watch"
)

// Handler handles requests
type Handler struct {
	TemplateDirectory string
	LogRepository     *repository.LogRepository
	Watcher           *watch.Watcher
}
