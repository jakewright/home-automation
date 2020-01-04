package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.scene/domain"
	sceneproto "github.com/jakewright/home-automation/service.scene/proto"
)

// HandleDeleteScene deletes a scene and associated actions
func HandleDeleteScene(body *sceneproto.DeleteSceneRequest) (*sceneproto.DeleteSceneResponse, error) {
	if body.SceneId == 0 {
		return nil, errors.BadRequest("scene_id empty")
	}

	// Find associated actions
	var actions []*domain.Action
	where := map[string]interface{}{
		"scene_id": body.SceneId,
	}
	if err := database.Find(&actions, where); err != nil {
		return nil, err
	}

	// Delete the actions
	if len(actions) > 0 {
		actionIDs := make([]uint, len(actions))
		for i, a := range actions {
			actionIDs[i] = a.ID
		}

		if err := database.Delete(&domain.Action{}, actionIDs); err != nil {
			return nil, err
		}
	}

	// Delete the scene
	if err := database.Delete(&domain.Scene{}, body.SceneId); err != nil {
		return nil, err
	}

	slog.Info("Deleted scene %d", body.SceneId)

	return &sceneproto.DeleteSceneResponse{}, nil
}
