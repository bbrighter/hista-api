package symptoms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateConditionEventAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var params = ConditionEventRequestParams{Date: time.Now()}
	resp, err := service.CreateConditionEvent(ctx, params)

	assert.NoError(t, err)
	assert.NotEqualValues(t, 0, resp.ID)
}

func TestGetConditionEventsAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)
	service.testCreateConditionEvent(t)

	events, _ := service.GetConditionEvents(ctx)
	assert.Len(t, events.ConditionEvents, 1)
}

func TestGetConditionEventAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)
	service.testCreateConditionEvent(t)

	var event ConditionEventResponse
	var err error
	event, err = service.GetConditionEvent(ctx, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, event.ID, 1)

	// Not found
	_, err = service.GetConditionEvent(ctx, 1000)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}

func TestPatchDateAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)
	service.testCreateConditionEvent(t)

	var err error
	var params = ConditionEventRequestParams{Date: time.Now()}
	err = service.PatchDate(ctx, 1, params)
	assert.NoError(t, err)
}

func TestDeleteConditionEventAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)
	service.testCreateConditionEvent(t)

	var err error
	err = service.DeleteConditionEvent(ctx, 1)
	assert.NoError(t, err)

	// Not found
	err = service.DeleteConditionEvent(ctx, 1000)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}
