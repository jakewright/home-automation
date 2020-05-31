package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/slog"
	scenedef "github.com/jakewright/home-automation/service.scene/def"
	"github.com/jakewright/home-automation/service.scene/domain"
)

// CreateScene persists a new scene
func (h *Handler) CreateScene(r *request, body *scenedef.CreateSceneRequest) (*scenedef.CreateSceneResponse, error) {
	actions := make([]*domain.Action, len(body.Actions))
	for i, a := range body.Actions {
		actions[i] = &domain.Action{
			Stage:          int(a.Stage),
			Sequence:       int(a.Sequence),
			Func:           a.Func,
			ControllerName: a.ControllerName,
			DeviceID:       a.DeviceId,
			Command:        a.Command,
			Property:       a.Property,
			PropertyValue:  a.PropertyValue,
			PropertyType:   a.PropertyType,
		}

		if err := actions[i].Validate(); err != nil {
			return nil, err
		}
	}

	scene := &domain.Scene{
		Name:    body.Name,
		Actions: actions,
	}

	if err := database.Create(scene); err != nil {
		return nil, err
	}

	slog.Infof("Created new scene %d", scene.ID)

	return &scenedef.CreateSceneResponse{
		Scene: scene.ToProto(),
	}, nil
}
