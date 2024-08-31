package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (service *Service) createTestSymptom(t *testing.T) func(t *testing.T) {
	id, cleanupEvent := service.createTestEvent(t)
	ctx := context.TODO()
	resp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "cat"})
	assert.NoError(t, err)
	conditionId, err := service.symtpoms.CreateConditionBySymptomName(id, "name", resp.ID)
	assert.NoError(t, err)

	cleanup := func(t *testing.T) {
		_, err := service.DeleteCondition(ctx, conditionId)
		assert.NoError(t, err)
		cleanupEvent(t)
	}

	return cleanup
}

func TestGetSymptoms(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, _ := service.GetSymptoms(ctx)
	assert.Len(t, resp.Categories, 0)

	cleanup := service.createTestSymptom(t)
	defer cleanup(t)

	resp, _ = service.GetSymptoms(ctx)
	assert.Len(t, resp.Categories, 1)
	assert.Equal(t, "cat", resp.Categories[0].Name)
	assert.Len(t, resp.Categories[0].Symptoms, 1)
	assert.Equal(t, "name", resp.Categories[0].Symptoms[0].Name)
}

func TestPostSymptomCategory(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "cat"})
	assert.NoError(t, err)
	assert.Greater(t, resp.ID, uint(0))

}
