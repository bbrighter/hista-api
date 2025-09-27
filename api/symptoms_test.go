package api

import (
	"context"
	"testing"

	entity "encore.app/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (service *Service) createTestSymptom(ctx context.Context, t *testing.T) (condId uint, symptomId uint, catId uint) {
	id := service.createTestEvent(ctx, t)
	catResp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "cat"})
	require.NoError(t, err)
	var name string = "name"
	resp, err := service.PostCondition(ctx, id, ConditionRequestParams{SymptomName: &name, CategoryID: &catResp.ID})
	require.NoError(t, err)

	return resp.Condition.ID, resp.Condition.Symptom.ID, resp.Condition.Symptom.CategoryID
}

func TestGetSymptoms(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, _ := service.GetSymptoms(ctx)
	assert.Len(t, resp.Categories, 0)

	service.createTestSymptom(ctx, t)

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

	_, symptomId, _ := service.createTestSymptom(ctx, t)

	err := service.PatchSymptomName(ctx, symptomId, PatchSymptomNameParams{Name: "new name"})
	assert.NoError(t, err)

	var symptom entity.Symptom
	service.DB.Take(&symptom, symptomId)
	assert.Equal(t, "new name", symptom.Name)
}

func TestPatchSymptomCategory(t *testing.T) {
	service, ctx := initAPITest(t)

	_, symptomId, _ := service.createTestSymptom(ctx, t)

	resp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "new cat"})
	assert.NoError(t, err)
	defer service.DeleteSymptomCategory(ctx, resp.ID)

	err = service.PatchSymptomCategory(ctx, symptomId, PatchSymptomCategoryParams{ToCategoryID: resp.ID})
	assert.NoError(t, err)

	var symptom entity.Symptom
	service.DB.Take(&symptom, symptomId)
	assert.Equal(t, resp.ID, symptom.SymptomCategoryID)
}

func TestPatchCategoryName(t *testing.T) {
	service, ctx := initAPITest(t)

	_, symptomId, catId := service.createTestSymptom(ctx, t)

	err := service.PatchCategoryName(ctx, symptomId, PatchCategoryNameParams{Name: "new cat name"})
	assert.NoError(t, err)

	var cat entity.SymptomCategory
	service.DB.Take(&cat, catId)
	assert.Equal(t, "new cat name", cat.Name)
}

func TestDeleteSymptomCategory(t *testing.T) {
	service, ctx := initAPITest(t)

	_, _, catId := service.createTestSymptom(ctx, t)

	err := service.DeleteSymptomCategory(ctx, catId)
	assert.Error(t, err)
	rows := service.DB.Take(&entity.SymptomCategories{}, catId).RowsAffected
	assert.EqualValues(t, 1, rows)

	resp, err := service.PostSymptomCategory(ctx, PostSymptomCategoryRequest{Name: "new cat"})
	assert.NoError(t, err)
	err = service.DeleteSymptomCategory(ctx, resp.ID)
	assert.NoError(t, err)
	rows = service.DB.Take(&entity.SymptomCategories{}, resp.ID).RowsAffected
	assert.EqualValues(t, 0, rows)
}
