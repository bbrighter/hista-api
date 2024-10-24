package api

import (
	"testing"
	"time"

	entity "encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateStatus(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.PostStatus(ctx, DateParam{time.Now()})
	defer service.DeleteStatus(ctx, resp.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(1))
}

func TestAddStatus(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	var resp, morningResp, eveningResp entity.StatusResponse
	resp, err = service.PostStatus(ctx, DateParam{time.Now()})
	defer service.DeleteStatus(ctx, resp.ID)
	assert.NoError(t, err)

	params := StatusParams{
		TimeOfDay: entity.Morning,
		Fitness:   entity.Good,
		Sleep:     entity.Good,
	}
	morningResp, err = service.PutStatus(ctx, resp.ID, params)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, morningResp.Morning.ID, uint(1))

	// Error if adding second time without fitting MorningID
	_, err = service.PutStatus(ctx, resp.ID, params)
	assert.Error(t, err)

	eveningParams := StatusParams{
		TimeOfDay: entity.Evening,
		Fitness:   entity.Middle,
	}
	eveningResp, err = service.PutStatus(ctx, resp.ID, eveningParams)
	assert.NoError(t, err)
	assert.Equal(t, entity.Middle, eveningResp.Evening.Fitness)
	assert.Equal(t, entity.Good, eveningResp.Morning.Fitness)

}

func TestAddStatusIdempotent(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.PostStatus(ctx, DateParam{time.Now()})
	defer service.DeleteStatus(ctx, resp.ID)
	assert.NoError(t, err)

	params := StatusParams{
		TimeOfDay: entity.Morning,
		Fitness:   entity.Good,
		Sleep:     entity.Good,
	}
	statusResp, err := service.PutStatus(ctx, resp.ID, params)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, statusResp.Morning.ID, uint(1))
	params.ID = statusResp.Morning.ID
	newStatusResp, err := service.PutStatus(ctx, resp.ID, params)
	assert.NoError(t, err)
	assert.Equal(t, newStatusResp.Morning.ID, statusResp.Morning.ID)
}

func TestListStatus(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, _ := service.ListStatus(ctx)
	assert.Len(t, resp.Statuses, 0)
}
