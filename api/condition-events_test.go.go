package api

import (
	"context"
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

var testEvent = new(entity.ConditionEvent)

func (service *Service) createTestEvent(t *testing.T) func(t *testing.T) {
	ctx := context.TODO()
	var params = ConditionEventRequestParams{Date: time.Now()}
	resp, err := service.CreateConditionEvent(ctx, params)
	testEvent.ID = resp.ID
	assert.NoError(t, err)
	cleanup := func(t *testing.T) {
		err = service.DeleteConditionEvent(ctx, resp.ID)
		assert.NoError(t, err)
		testEvent = new(entity.ConditionEvent)
	}
	return cleanup
}

func TestCreateConditionEvent(t *testing.T) {
	service, ctx := initAPITest(t)

	var params = ConditionEventRequestParams{Date: time.Now()}
	resp, err := service.CreateConditionEvent(ctx, params)
	defer service.DeleteConditionEvent(ctx, resp.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.ID)
}

func TestGetConditionEvents(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.GetConditionEvents(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.ConditionEvents, 0)

	cleanup := service.createTestEvent(t)
	defer cleanup(t)

	resp, err = service.GetConditionEvents(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.ConditionEvents, 1)
}

func TestGetConditionEvent(t *testing.T) {
	service, ctx := initAPITest(t)

	_, err := service.GetConditionEvent(ctx, 100)
	assert.EqualError(t, err, "not_found: not found")

	cleanup := service.createTestEvent(t)
	defer cleanup(t)

	resp, err := service.GetConditionEvent(ctx, testEvent.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, testEvent.ID, resp.ID)
	assert.True(t, time.Now().After(resp.Date))
}

func TestPatchConditionEvent(t *testing.T) {
	service, ctx := initAPITest(t)

	var params ConditionEventRequestParams
	err := service.PatchDate(ctx, 10, params)
	assert.EqualError(t, err, "not_found: not found")

	cleanup := service.createTestEvent(t)
	defer cleanup(t)
	err = service.PatchDate(ctx, testEvent.ID, params)
	assert.Error(t, err)

	params.Date = time.Now()
	err = service.PatchDate(ctx, testEvent.ID, params)
	assert.NoError(t, err)
}

func TestDeleteConditionEvent(t *testing.T) {
	service, ctx := initAPITest(t)

	err := service.DeleteConditionEvent(ctx, 10)
	assert.EqualError(t, err, "not_found: not found")

	cleanup := service.createTestEvent(t)
	defer cleanup(t)
	err = service.DeleteConditionEvent(ctx, testEvent.ID)
	assert.NoError(t, err)
}

func TestPostCondition(t *testing.T) {
	service, ctx := initAPITest(t)

	var params ConditionRequestParams
	var err error
	var resp entity.PostConditionResponse

	_, err = service.PostCondition(ctx, 10, params)
	assert.EqualError(t, err, "not_found: not found")

	cleanup := service.createTestEvent(t)
	defer cleanup(t)
	_, err = service.PostCondition(ctx, testEvent.ID, params)
	assert.Error(t, err)

	// Name + CategoryId
	var name string = "name"
	catResp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "cat"})
	assert.NoError(t, err)
	params.SymptomName = &name
	params.CategoryID = &catResp.ID

	resp, err = service.PostCondition(ctx, testEvent.ID, params)
	defer service.DeleteCondition(ctx, resp.Condition.ID)
	assert.NoError(t, err)

	// SymptomId
	params.SymptomName = nil
	params.CategoryID = nil
	params.SymptomID = &resp.Symptoms.Categories[0].Symptoms[0].ID

	resp, err = service.PostCondition(ctx, testEvent.ID, params)
	defer service.DeleteCondition(ctx, resp.Condition.ID)
	assert.NoError(t, err)
}
