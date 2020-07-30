package consumer

import (
	"testing"

	"gotest.tools/assert"

	"github.com/jakewright/home-automation/services/scene/domain"
)

func Test_constructStages(t *testing.T) {
	scene := &domain.Scene{
		Actions: []*domain.Action{
			{Stage: 1, Sequence: 1},
			{Stage: 1, Sequence: 2},
			{Stage: 1, Sequence: 3},
			{Stage: 2, Sequence: 1},
			{Stage: 3, Sequence: 1},
			{Stage: 3, Sequence: 1},
			{Stage: 6, Sequence: 2},
		},
	}

	stages := constructStages(scene)

	assert.Equal(t, 4, len(stages))
	assert.Equal(t, 3, len(stages[0]))
	assert.Equal(t, 1, len(stages[1]))
	assert.Equal(t, 2, len(stages[2]))
	assert.Equal(t, 1, len(stages[3]))
}
