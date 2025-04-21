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
	catResp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "cat"})
	assert.NoError(t, err)
	var name string = "name"
	condition, cats, err := service.conditions.Create(testEvent.ID, &name, nil, &catResp.ID)
	testCondition = &condition
	testSymptom.ID = cats[0].Symptoms[0].ID
	testSymptom.SymptomCategoryID = catResp.ID
	assert.NoError(t, err)

	cleanup := func(t *testing.T) {
		_, err := service.DeleteCondition(ctx, condition.ID)
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

func TestPatchSymptomName(t *testing.T) {
	service, ctx := initAPITest(t)

	cleanup := service.createTestSymptom(t)
	defer cleanup(t)

	err := service.PatchSymptomName(ctx, testSymptom.ID, PatchSymptomNameParams{Name: "new name"})
	assert.NoError(t, err)

	var symptom entity.Symptom
	service.DB.Take(&symptom, testSymptom.ID)
	assert.Equal(t, "new name", symptom.Name)
}

func TestPatchSymptomCategory(t *testing.T) {
	service, ctx := initAPITest(t)

	cleanup := service.createTestSymptom(t)
	defer cleanup(t)

	resp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "new cat"})
	assert.NoError(t, err)
	defer service.DeleteSymptomCategory(ctx, resp.ID)

	err = service.PatchSymptomCategory(ctx, testSymptom.ID, PatchSymptomCategoryParams{ToCategoryID: resp.ID})
	assert.NoError(t, err)

	var symptom entity.Symptom
	service.DB.Take(&symptom, testSymptom.ID)
	assert.Equal(t, resp.ID, symptom.SymptomCategoryID)
}

func TestPatchCategoryName(t *testing.T) {
	service, ctx := initAPITest(t)

	cleanup := service.createTestSymptom(t)
	defer cleanup(t)

	err := service.PatchCategoryName(ctx, testSymptom.SymptomCategoryID, PatchCategoryNameParams{Name: "new cat name"})
	assert.NoError(t, err)

	var cat entity.SymptomCategory
	service.DB.Debug().Take(&cat, testSymptom.SymptomCategoryID)
	assert.Equal(t, "new cat name", cat.Name)
}

func TestDeleteSymptomCategory(t *testing.T) {
	service, ctx := initAPITest(t)

	cleanup := service.createTestSymptom(t)
	defer cleanup(t)

	err := service.DeleteSymptomCategory(ctx, testSymptom.SymptomCategoryID)
	assert.Error(t, err)
	rows := service.DB.Take(&entity.SymptomCategories{}, testSymptom.SymptomCategoryID).RowsAffected
	assert.EqualValues(t, 1, rows)

	resp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "new cat"})
	assert.NoError(t, err)
	err = service.DeleteSymptomCategory(ctx, resp.ID)
	assert.NoError(t, err)
	rows = service.DB.Take(&entity.SymptomCategories{}, resp.ID).RowsAffected
	assert.EqualValues(t, 0, rows)
}
