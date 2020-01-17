package consumer

import (
	"encoding/json"
	"sort"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/dsync"
	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.scene/domain"
)

// HandleSetSceneEvent sets the scene
func HandleSetSceneEvent(event *firehose.Event) firehose.Result {
	var body struct {
		SceneID uint32 `json:"scene_id"`
	}
	if err := json.Unmarshal(event.Payload, &body); err != nil {
		return firehose.Discard(errors.WithMessage(err, "failed to unmarshal payload"))
	}

	metadata := map[string]string{
		"scene_id": strconv.Itoa(int(body.SceneID)),
	}

	scene := &domain.Scene{}
	if err := database.Find(&scene, body.SceneID); err != nil {
		return firehose.Discard(errors.WithMetadata(err, metadata))
	}

	if scene == nil {
		err := errors.NotFound("scene not found", metadata)
		return firehose.Discard(err)
	}

	lock, err := dsync.Lock("scene", body.SceneID)
	if err != nil {
		return firehose.Fail(errors.WithMetadata(err, metadata))
	}
	defer lock.Unlock()

	slog.Infof("Setting scene %d...", body.SceneID)

	// Organise the actions into stages
	stages := constructStages(scene)

	for _, stage := range stages {
		var g errgroup.Group
		for _, action := range stage {
			g.Go(action.Perform)
		}
		if err := g.Wait(); err != nil {
			return firehose.Fail(err)
		}
	}

	return firehose.Success()
}

func constructStages(scene *domain.Scene) [][]*domain.Action {
	m := make(map[int][]*domain.Action)

	// Sort the actions into buckets
	for _, action := range scene.Actions {
		m[action.Stage] = append(m[action.Stage], action)
	}

	// Turn the map into a slice of slices
	var stages [][]*domain.Action
	for _, actions := range m {
		stages = append(stages, actions)
	}

	// Sort the slices by stage number
	sort.Slice(stages, func(i, j int) bool {
		return stages[i][0].Stage < stages[j][0].Stage
	})

	return stages
}
