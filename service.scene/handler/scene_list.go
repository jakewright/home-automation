package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/service.scene/domain"
	"github.com/jakewright/home-automation/service.scene/external"
)

// HandleListScenes lists all scenes in the external
func HandleListScenes(req *external.ListScenesRequest) (*external.ListScenesResponse, error) {
	where := make(map[string]interface{})
	if req.OwnerId > 0 {
		where["owner_id"] = req.OwnerId
	}

	var scenes []*domain.Scene
	if err := database.Find(&scenes, where); err != nil {
		return nil, err
	}

	protos := make([]external.Scene, len(scenes))
	for i, s := range scenes {
		protos[i] = s.ToProto()
	}

	return &external.ListScenesResponse{
		Scenes: protos,
	}, nil
}
