package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.scene/domain"
	"github.com/jakewright/home-automation/service.scene/external"
)

// HandleCreateScene persists a new scene
func HandleCreateScene(req *external.CreateSceneRequest) (*external.CreateSceneResponse, error) {
	actions := make([]*domain.Action, len(req.Actions))
	for i, a := range req.Actions {
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
		Name:    req.Name,
		Actions: actions,
	}

	if err := database.Create(scene); err != nil {
		return nil, err
	}

	slog.Infof("Created new scene %d", scene.ID)

	return &external.CreateSceneResponse{
		Scene: scene.ToProto(),
	}, nil
}
