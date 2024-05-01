package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMealsAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	resp, err := service.GetMeals(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Meals, 0)
}

func TestPostMealAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var params = MealParams{Date: time.Now()}
	resp, err := service.PostMeal(ctx, params)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(1))
}

func TestGetMealAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var err error
	_, err = service.GetMeal(ctx, 1)
	assert.Error(t, err)

	var id uint = service.testCreateMeal(t).ID
	var resp MealResponse
	resp, err = service.GetMeal(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, resp.ID)
}

func TestDeleteMealAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)

	var err error
	err = service.DeleteMeal(ctx, 1)
	assert.Error(t, err)

	var id uint = service.testCreateMeal(t).ID
	err = service.DeleteMeal(ctx, id)
	assert.NoError(t, err)
}

func TestPatchMealTimeAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	defer teardown(t)
	var id uint = service.testCreateMeal(t).ID

	var params = MealParams{Date: time.Now()}
	var err error = service.PatchMealTime(ctx, id, params)
	assert.NoError(t, err)
}
