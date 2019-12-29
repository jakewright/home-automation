package domain

import "github.com/jinzhu/gorm"

type Scene struct {
	gorm.Model
	Name string
	Actions []Action
}

type Action struct {
	gorm.Model
	SceneID int
	Stage int
	Index int

	Function string
	ControllerName string
	Command string
	Property string
	PropertyValue string
}
