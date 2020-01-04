package domain

import (
	"github.com/jinzhu/gorm"

	"github.com/jakewright/home-automation/libraries/go/util"
	sceneproto "github.com/jakewright/home-automation/service.scene/proto"
)

// Scene represents a set of actions
type Scene struct {
	gorm.Model
	Name    string
	Actions []*Action
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

// ToProto marshals to the proto type
func (s *Scene) ToProto() *sceneproto.Scene {
	actions := make([]*sceneproto.Action, len(s.Actions))
	for i, a := range s.Actions {
		actions[i] = a.ToProto()
	}

	return &sceneproto.Scene{
		Id:        uint32(s.ID),
		Name:      s.Name,
		Actions:   actions,
		CreatedAt: util.TimeToProto(s.CreatedAt),
		UpdatedAt: util.TimeToProto(s.UpdatedAt),
		DeletedAt: util.PTimeToProto(s.DeletedAt),
	}
}

// ToProto marshals to the proto type
func (a *Action) ToProto() *sceneproto.Action {
	return &sceneproto.Action{
		Id:             uint32(a.ID),
		Stage:          int32(a.Stage),
		Sequence:       int32(a.Sequence),
		Func:           a.Func,
		ControllerName: a.ControllerName,
		Command:        a.Command,
		Property:       a.Property,
		PropertyValue:  a.PropertyValue,
		CreatedAt:      util.TimeToProto(a.CreatedAt),
		UpdatedAt:      util.TimeToProto(a.UpdatedAt),
		DeletedAt:      util.PTimeToProto(a.DeletedAt),
	}
}
