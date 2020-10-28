package routes

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/oops"
	scenedef "github.com/jakewright/home-automation/services/scene/def"
	"github.com/jakewright/home-automation/services/scene/domain"
)

// ReadScene returns the scene with the given ID
func (c *Controller) ReadScene(ctx context.Context, body *scenedef.ReadSceneRequest) (*scenedef.ReadSceneResponse, error) {
	scene := &domain.Scene{}
	if err := c.Database.Find(&scene, body.SceneId); err != nil {
		return nil, oops.WithMessage(err, "failed to find")
	}

	return &scenedef.ReadSceneResponse{
		Scene: scene.ToProto(),
	}, nil
}
