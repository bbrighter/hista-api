package hista

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"github.com/stretchr/testify/assert"
)

func (service *Service) createTestMeal(ctx context.Context, t *testing.T) (uint, func()) {
	meal, err := service.PostMeal(ctx, TEST_PIID, entity.MealParams{})
	assert.NoError(t, err)
	cleanUp := func() {
		service.DeleteMeal(ctx, TEST_PIID, meal.ID)
	}
	return meal.ID, cleanUp
}

func TestGetMealsAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.ListMeals(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, resp.Meals, 0)

	_, cleanup := service.createTestMeal(ctx, t)
	defer cleanup()

	resp, err = service.ListMeals(ctx, TEST_PIID)
	assert.NoError(t, err)
	assert.Len(t, resp.Meals, 1)
}

func TestPostMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var now = time.Now()
	var params = entity.MealParams{Date: &now}
	resp, err := service.PostMeal(ctx, TEST_PIID, params)
	defer service.DeleteMeal(ctx, TEST_PIID, resp.ID)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.ID, uint(1))
}

func TestGetMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	_, err = service.GetMeal(ctx, TEST_PIID, 100)
	assert.Error(t, err)

	id, cleanup := service.createTestMeal(ctx, t)
	defer cleanup()

	resp, err := service.GetMeal(ctx, TEST_PIID, id)
	assert.NoError(t, err, err)
	assert.EqualValues(t, id, resp.ID)
}

func TestDeleteMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	_, err = service.DeleteMeal(ctx, TEST_PIID, 10000)
	assert.Error(t, err)

	id, cleanup := service.createTestMeal(ctx, t)
	defer cleanup()
	ings, err := service.DeleteMeal(ctx, TEST_PIID, id)
	assert.NoError(t, err)
	assert.Len(t, ings.Ingredients, 0)
}

func TestPatchMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var params entity.MealParams
	var now time.Time = time.Now()
	params.Date = &now

	err := service.PatchMeal(ctx, TEST_PIID, 10000, params)
	assert.Error(t, err)
	id, cleanup := service.createTestMeal(ctx, t)
	defer cleanup()

	err = service.PatchMeal(ctx, TEST_PIID, id, params)
	assert.NoError(t, err)

	var stressLevel uint8 = 2
	params.StressLevel = &stressLevel
	err = service.PatchMeal(ctx, TEST_PIID, id, params)
	assert.NoError(t, err)

}

func TestGetFoods(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.GetFoods(ctx, TEST_PIID, 1)
	assert.NoError(t, err)

	id, cleanup := service.createTestMeal(ctx, t)
	defer cleanup()

	resp, err = service.GetFoods(ctx, TEST_PIID, id)
	assert.NoError(t, err)
	assert.Len(t, resp.Foods, 0)
}

func TestPostFood(t *testing.T) {
	service, ctx := initAPITest(t)

	mealId, cleanup := service.createTestMeal(ctx, t)
	defer cleanup()

	var params = FoodParams{IngredientName: "New", Condition: entity.Cooked}
	food, err := service.PostFood(ctx, TEST_PIID, mealId, params)
	defer service.DeleteFood(ctx, TEST_PIID, food.Food.ID)
	assert.NoError(t, err)
}

func TestDeleteFoodAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	_, err := service.DeleteFood(ctx, TEST_PIID, 100)
	assert.EqualError(t, err, "not_found: not found")

	mealId, cleanup := service.createTestMeal(ctx, t)
	defer cleanup()
	var params = FoodParams{IngredientName: "New", Condition: entity.Cooked}
	food, _ := service.PostFood(ctx, TEST_PIID, mealId, params)

	ing, err := service.DeleteFood(ctx, TEST_PIID, food.Food.ID)
	assert.NoError(t, err)
	assert.Len(t, ing.Ingredients, 0)
}
