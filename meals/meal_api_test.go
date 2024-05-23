package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMealsAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.GetMeals(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Meals, 1)
}

func TestPostMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var params = MealParams{Date: time.Now()}
	resp, err := service.PostMeal(ctx, params)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(1))

	// Cleanup
	err = service.deleteMeal(resp.ID)
	assert.NoError(t, err)
}

func TestGetMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	_, err = service.GetMeal(ctx, 100)
	assert.Error(t, err)

	var resp MealResponse
	resp, err = service.GetMeal(ctx, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, resp.ID)
}

func TestDeleteMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	err = service.DeleteMeal(ctx, 100)
	assert.Error(t, err)

	var id uint
	id, err = service.createMeal(time.Now())
	assert.NoError(t, err)
	err = service.DeleteMeal(ctx, id)
	assert.NoError(t, err)
}

func TestPatchMealTimeAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var params = MealParams{Date: time.Now()}
	var err error = service.PatchMealTime(ctx, 1, params)
	assert.NoError(t, err)
}
