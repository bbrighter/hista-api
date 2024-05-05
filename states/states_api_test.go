package states

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateStateAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var params = StateRequestParams{Date: time.Now()}
	resp, err := service.CreateState(ctx, params)

	assert.NoError(t, err)
	assert.NotEqualValues(t, 0, resp.ID)
}

func TestGetStatesAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)
	service.testCreateState(t)

	states, _ := service.GetStates(ctx)
	assert.Len(t, states.States, 1)
}

func TestGetStateAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)
	service.testCreateState(t)

	var state StateResponse
	var err error
	state, err = service.GetState(ctx, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, state.ID, 1)

	// Not found
	_, err = service.GetState(ctx, 1000)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}

func TestDeleteStateAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)
	service.testCreateState(t)

	var err error
	err = service.DeleteState(ctx, 1)
	assert.NoError(t, err)

	// Not found
	err = service.DeleteState(ctx, 1000)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}
