package symptoms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateConditionEventAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var params = ConditionEventRequestParams{Date: time.Now()}
	resp, err := service.CreateConditionEvent(ctx, params)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(2))

	// clean up
	var event ConditionEvent
	service.db.Last(&event)
	err = event.delete(service)
	assert.NoError(t, err)
}

func TestGetConditionEventsAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	events, _ := service.GetConditionEvents(ctx)
	assert.GreaterOrEqual(t, len(events.ConditionEvents), 1)
}

func TestGetConditionEventAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var event ConditionEventResponse
	var err error
	event, err = service.GetConditionEvent(ctx, testEvent.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, event.ID, testEvent.ID)

	// Not found
	_, err = service.GetConditionEvent(ctx, 1000)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}

func TestPatchDateAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	var params = ConditionEventRequestParams{Date: time.Now()}
	err = service.PatchDate(ctx, testEvent.ID, params)
	assert.NoError(t, err)

	// Cleanup
	params.Date = testEvent.Date
	service.PatchDate(ctx, testEvent.ID, params)
}

func TestDeleteConditionEventAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	err = service.DeleteConditionEvent(ctx, testEvent.ID)
	assert.NoError(t, err)

	// Not found
	err = service.DeleteConditionEvent(ctx, 1000)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not found")
}
