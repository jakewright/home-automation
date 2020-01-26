package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	scenedef "github.com/jakewright/home-automation/service.scene/def"
	"github.com/jakewright/home-automation/service.scene/domain"
)

// HandleListScenes lists all scenes in the database
func HandleListScenes(req *scenedef.ListScenesRequest) (*scenedef.ListScenesResponse, error) {
	where := make(map[string]interface{})
	if req.OwnerId > 0 {
		where["owner_id"] = req.OwnerId
	}

	var scenes []*domain.Scene
	if err := database.Find(&scenes, where); err != nil {
		return nil, err
	}

	protos := make([]scenedef.Scene, len(scenes))
	for i, s := range scenes {
		protos[i] = s.ToProto()
	}

	return &scenedef.ListScenesResponse{
		Scenes: protos,
	}, nil
}
