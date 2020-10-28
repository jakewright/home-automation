package routes

import (
	"context"

	scenedef "github.com/jakewright/home-automation/services/scene/def"
	"github.com/jakewright/home-automation/services/scene/domain"
)

// ListScenes lists all scenes in the database
func (c *Controller) ListScenes(ctx context.Context, body *scenedef.ListScenesRequest) (*scenedef.ListScenesResponse, error) {
	where := make(map[string]interface{})
	if body.OwnerId > 0 {
		where["owner_id"] = body.OwnerId
	}

	var scenes []*domain.Scene
	if err := c.Database.Find(&scenes, where); err != nil {
		return nil, err
	}

	protos := make([]*scenedef.Scene, len(scenes))
	for i, s := range scenes {
		protos[i] = s.ToProto()
	}

	return &scenedef.ListScenesResponse{
		Scenes: protos,
	}, nil
}
