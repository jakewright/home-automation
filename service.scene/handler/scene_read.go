package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/oops"
	scenedef "github.com/jakewright/home-automation/service.scene/def"
	"github.com/jakewright/home-automation/service.scene/domain"
)

// ReadScene returns the scene with the given ID
func (h *Handler) ReadScene(r *request, body *scenedef.ReadSceneRequest) (*scenedef.ReadSceneResponse, error) {
	scene := &domain.Scene{}
	if err := database.Find(&scene, body.SceneId); err != nil {
		return nil, oops.WithMessage(err, "failed to find")
	}

	return &scenedef.ReadSceneResponse{
		Scene: scene.ToProto(),
	}, nil
}
