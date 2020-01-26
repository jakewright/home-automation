package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/service.scene/domain"
	"github.com/jakewright/home-automation/service.scene/external"
)

// HandleSetScene emits an event to trigger the scene to be set asynchronously
func HandleSetScene(body *external.SetSceneRequest) (*external.SetSceneResponse, error) {
	scene := &domain.Scene{}
	if err := database.Find(&scene, body.SceneId); err != nil {
		return nil, err
	}

	if scene == nil {
		return nil, errors.NotFound("Scene not found")
	}

	if err := (&external.SetSceneEvent{
		SceneId: body.SceneId,
	}).Publish(); err != nil {
		return nil, err
	}

	return &external.SetSceneResponse{}, nil
}
