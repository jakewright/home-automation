package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.scene/domain"
	"github.com/jakewright/home-automation/service.scene/external"
)

// HandleDeleteScene deletes a scene and associated actions
func HandleDeleteScene(req *external.DeleteSceneRequest) (*external.DeleteSceneResponse, error) {
	if req.SceneId == 0 {
		return nil, errors.BadRequest("scene_id empty")
	}

	// Delete the scene
	if err := database.Delete(&domain.Scene{}, req.SceneId); err != nil {
		return nil, err
	}

	slog.Infof("Deleted scene %d", req.SceneId)
	return &external.DeleteSceneResponse{}, nil
}
