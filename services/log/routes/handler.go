package routes

import (
	"github.com/jakewright/home-automation/services/log/repository"
	"github.com/jakewright/home-automation/services/log/watch"
)

// Handler handles requests
type Handler struct {
	TemplateDirectory string
	LogRepository     *repository.LogRepository
	Watcher           *watch.Watcher
}
