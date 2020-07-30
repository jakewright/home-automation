package domain

import (
	"time"

	scenedef "github.com/jakewright/home-automation/services/scene/def"
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
func (s *Scene) ToProto() *scenedef.Scene {
	actions := make([]*scenedef.Action, len(s.Actions))
	for i, a := range s.Actions {
		actions[i] = a.ToProto()
	}

	return &scenedef.Scene{
		Id:        s.ID,
		Name:      s.Name,
		OwnerId:   s.OwnerID,
		Actions:   actions,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
