package routes

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
	scenedef "github.com/jakewright/home-automation/services/scene/def"
	"github.com/jakewright/home-automation/services/scene/domain"
)

// DeleteScene deletes a scene and associated actions
func (c *Controller) DeleteScene(ctx context.Context, body *scenedef.DeleteSceneRequest) (*scenedef.DeleteSceneResponse, error) {
	if body.SceneId == 0 {
		return nil, oops.BadRequest("scene_id empty")
	}

	// Delete the scene
	if err := c.Database.Delete(&domain.Scene{}, body.SceneId); err != nil {
		return nil, err
	}

	slog.Infof("Deleted scene %d", body.SceneId)
	return &scenedef.DeleteSceneResponse{}, nil
}
