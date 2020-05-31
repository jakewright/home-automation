package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
	scenedef "github.com/jakewright/home-automation/service.scene/def"
	"github.com/jakewright/home-automation/service.scene/domain"
)

// DeleteScene deletes a scene and associated actions
func (h *Handler) DeleteScene(r *request, body *scenedef.DeleteSceneRequest) (*scenedef.DeleteSceneResponse, error) {
	if body.SceneId == 0 {
		return nil, oops.BadRequest("scene_id empty")
	}

	// Delete the scene
	if err := database.Delete(&domain.Scene{}, body.SceneId); err != nil {
		return nil, err
	}

	slog.Infof("Deleted scene %d", body.SceneId)
	return &scenedef.DeleteSceneResponse{}, nil
}
