package api

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateStatus(t *testing.T) {
	service, ctx := initAPITest(t)

	params := StatusParams{
		TimeOfDay: entity.Morning,
		Date:      time.Now(),
		Fitness:   entity.Good,
		Sleep:     entity.Good,
	}

	resp, err := service.CreateStatus(ctx, params)
	defer service.DeleteStatus(ctx, resp.ID)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(1))
}

func TestCreateStatusEvening(t *testing.T) {
	service, ctx := initAPITest(t)

	params := StatusParams{
		TimeOfDay: entity.Evening,
		Date:      time.Now(),
		Fitness:   entity.Good,
	}

	resp, err := service.CreateStatus(ctx, params)
	defer service.DeleteStatus(ctx, resp.ID)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(1))
}

func TestListStatus(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, _ := service.ListStatus(ctx)
	assert.Len(t, resp.Statuses, 0)
}
