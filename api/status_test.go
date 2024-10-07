package api

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateStatus(t *testing.T) {
	service, ctx := initAPITest(t)

	// sleep := entity.Bad
	params := StatusParams{
		Date:      time.Now(),
		TimeOfDay: entity.Evening,
		Fitness:   entity.Good,
		Sleep:     nil,
	}

	resp, err := service.CreateStatus(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, entity.Good, resp.Fitness)
}

func TestListStatus(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, _ := service.ListStatus(ctx)
	assert.Len(t, resp.Statuses, 0)
}
