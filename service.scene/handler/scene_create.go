package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.scene/domain"
	sceneproto "github.com/jakewright/home-automation/service.scene/proto"
)

// HandleCreateScene persists a new scene
func HandleCreateScene(req *sceneproto.CreateSceneRequest) (*sceneproto.CreateSceneResponse, error) {
	actions := make([]*domain.Action, len(req.Actions))
	for i, a := range req.Actions {
		actions[i] = &domain.Action{
			Stage:          int(a.Stage),
			Sequence:       int(a.Sequence),
			Func:           a.Func,
			ControllerName: a.ControllerName,
			Command:        a.Command,
			Property:       a.PropertyValue,
			PropertyValue:  a.PropertyValue,
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

	return &sceneproto.CreateSceneResponse{
		Scene: scene.ToProto(),
	}, nil
}
