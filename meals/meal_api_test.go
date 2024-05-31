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

	var now = time.Now()
	var params = MealParams{Date: &now}
	resp, err := service.PostMeal(ctx, params)
	defer service.DeleteMeal(ctx, resp.ID)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(1))
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

	err = service.DeleteMeal(ctx, testMeal.ID)
	assert.NoError(t, err)
}

func TestPatchMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var meal = Meal{ID: testMeal.ID}
	defer service.DeleteMeal(ctx, testMeal.ID)

	var params MealParams
	var now time.Time = time.Now()
	params.Date = &now
	var err error
	err = service.PatchMeal(ctx, meal.ID, params)
	assert.NoError(t, err)

	var stressLevel uint8 = 2
	params.StressLevel = &stressLevel
	err = service.PatchMeal(ctx, meal.ID, params)
	assert.NoError(t, err)

	err = service.PatchMeal(ctx, 10000, params)
	assert.Error(t, err)
}
