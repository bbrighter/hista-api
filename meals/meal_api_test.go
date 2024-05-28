package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMealsAPI(t *testing.T) {
	service, ctx := initAPITest(t)
	service.initData()

	resp, err := service.GetMeals(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(resp.Meals), 1)
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
	resp, err = service.GetMeal(ctx, testMeal.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, testMeal.ID, resp.ID)
}

func TestDeleteMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	err = service.DeleteMeal(ctx, 10000)
	assert.Error(t, err)

	var id uint
	id, err = service.createMeal(time.Now())
	assert.NoError(t, err)
	err = service.DeleteMeal(ctx, id)
	assert.NoError(t, err)
}

func TestPatchMealTimeAPI(t *testing.T) {
	service, ctx := initAPITest(t)
	service.initData()

	var id uint
	var err error
	id, err = service.createMeal(time.Now())
	var params = MealParams{Date: time.Now()}
	service.PatchMealTime(ctx, id, params)
	assert.NoError(t, err)

	service.deleteMeal(id)
}
