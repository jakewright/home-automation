package domain

import (
	"time"

	"github.com/jakewright/home-automation/service.scene/external"
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
func (s *Scene) ToProto() external.Scene {
	actions := make([]external.Action, len(s.Actions))
	for i, a := range s.Actions {
		actions[i] = a.ToProto()
	}

	return external.Scene{
		Id:        s.ID,
		Name:      s.Name,
		OwnerId:   s.OwnerID,
		Actions:   actions,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
