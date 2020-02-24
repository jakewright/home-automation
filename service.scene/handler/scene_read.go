package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/errors"
	scenedef "github.com/jakewright/home-automation/service.scene/def"
	"github.com/jakewright/home-automation/service.scene/domain"
)

// HandleReadScene returns the scene with the given ID
func HandleReadScene(r *Request, body *scenedef.ReadSceneRequest) (*scenedef.ReadSceneResponse, error) {
	scene := &domain.Scene{}
	if err := database.Find(&scene, body.SceneId); err != nil {
		return nil, errors.WithMessage(err, "failed to find")
	}

	return &scenedef.ReadSceneResponse{
		Scene: scene.ToProto(),
	}, nil
}
