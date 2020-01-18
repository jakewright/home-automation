package domain

import (
	"time"

	"github.com/jakewright/home-automation/libraries/go/util"
	sceneproto "github.com/jakewright/home-automation/service.scene/proto"
)

// Scene represents a set of actions
type Scene struct {
	ID        uint32
	Name      string
	OwnerID   uint32
	Actions   []*Action
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToProto marshals to the proto type
func (s *Scene) ToProto() *sceneproto.Scene {
	actions := make([]*sceneproto.Action, len(s.Actions))
	for i, a := range s.Actions {
		actions[i] = a.ToProto()
	}

	return &sceneproto.Scene{
		Id:        s.ID,
		Name:      s.Name,
		OwnerId:   s.OwnerID,
		Actions:   actions,
		CreatedAt: util.TimeToProto(s.CreatedAt),
		UpdatedAt: util.TimeToProto(s.UpdatedAt),
	}
}
