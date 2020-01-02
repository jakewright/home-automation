package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/service.scene/domain"
	sceneproto "github.com/jakewright/home-automation/service.scene/proto"
)

// HandleCreateScene creates a new scene in the database
func HandleCreateScene(body *sceneproto.CreateSceneRequest) (*sceneproto.CreateSceneResponse, error) {
	scene := &domain.Scene{
		Name:    body.Name,
		Actions: nil,
	}

	if err := database.Create(scene).Error; err != nil {
		return nil, err
	}

	return &sceneproto.CreateSceneResponse{}, nil
}
