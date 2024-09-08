package api

import (
	"context"
	"testing"

	entity "encore.app/entity"
	"github.com/stretchr/testify/assert"
)

var testSymptom = new(entity.Symptom)
var testCondition = new(entity.Condition)

func (service *Service) createTestSymptom(t *testing.T) func(t *testing.T) {
	cleanupEvent := service.createTestEvent(t)
	ctx := context.TODO()
	resp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "cat"})
	assert.NoError(t, err)
	var name string = "name"
	conditon, cats, err := service.conditions.Create(testEvent.ID, &name, nil, &resp.ID)
	testCondition = &conditon
	testSymptom.ID = cats[0].Symptoms[0].ID
	assert.NoError(t, err)

	cleanup := func(t *testing.T) {
		_, err := service.DeleteCondition(ctx, conditon.ID)
		assert.NoError(t, err)
		cleanupEvent(t)
		testCondition = new(entity.Condition)
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
