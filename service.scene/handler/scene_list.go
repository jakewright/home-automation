package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/service.scene/domain"
	sceneproto "github.com/jakewright/home-automation/service.scene/proto"
)

// HandleListScenes lists all scenes in the database
func HandleListScenes(req *sceneproto.ListScenesRequest) (*sceneproto.ListScenesResponse, error) {
	where := make(map[string]interface{})
	if req.OwnerId > 0 {
		where["owner_id"] = req.OwnerId
	}

	var scenes []*domain.Scene
	if err := database.Find(&scenes, where); err != nil {
		return nil, err
	}

	protos := make([]*sceneproto.Scene, len(scenes))
	for i, s := range scenes {
		protos[i] = s.ToProto()
	}

	return &sceneproto.ListScenesResponse{
		Scenes: protos,
	}, nil
}
