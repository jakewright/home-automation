package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.scene/handler"
	sceneproto "github.com/jakewright/home-automation/service.scene/proto"
)

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.scene",
		Database:    true,
	})
	if err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	r := sceneproto.NewRouter()
	r.CreateScene = handler.HandleCreateScene
	r.ReadScene = handler.HandleReadScene
	r.ListScenes = handler.HandleListScenes
	r.DeleteScene = handler.HandleDeleteScene

	svc.Run(r)
}
