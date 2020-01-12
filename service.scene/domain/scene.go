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
