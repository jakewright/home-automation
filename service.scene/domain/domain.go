package domain

import "github.com/jinzhu/gorm"

// Scene represents a set of actions
type Scene struct {
	gorm.Model
	Name    string
	Actions []Action
}

// Action is a single step in a scene
type Action struct {
	gorm.Model
	SceneID  int
	Stage    int
	Sequence int

	Func           string
	ControllerName string
	Command        string
	Property       string
	PropertyValue  string
}
