package api

import (
	"context"
	"testing"
	"time"

	entity "encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func (service *Service) createTestMeal(t *testing.T) (uint, func()) {
	ctx := context.TODO()
	meal, err := service.PostMeal(ctx, entity.MealParams{})
	assert.NoError(t, err)
	cleanUp := func() {
		service.DeleteMeal(ctx, meal.ID)
	}
	return meal.ID, cleanUp
}

func TestGetMealsAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.GetMeals(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Meals, 0)

	_, cleanup := service.createTestMeal(t)
	defer cleanup()

	resp, err = service.GetMeals(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp.Meals, 1)
}

func TestPostMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var now = time.Now()
	var params = entity.MealParams{Date: &now}
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

	id, cleanup := service.createTestMeal(t)
	defer cleanup()

	resp, err := service.GetMeal(ctx, id)
	assert.NoError(t, err, err)
	assert.EqualValues(t, id, resp.ID)
}

func TestDeleteMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var err error
	err = service.DeleteMeal(ctx, 10000)
	assert.Error(t, err)

	id, cleanup := service.createTestMeal(t)
	defer cleanup()
	err = service.DeleteMeal(ctx, id)
	assert.NoError(t, err)
}

func TestPatchMealAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var params entity.MealParams
	var now time.Time = time.Now()
	params.Date = &now

	err := service.PatchMeal(ctx, 10000, params)
	assert.Error(t, err)
	id, cleanup := service.createTestMeal(t)
	defer cleanup()

	err = service.PatchMeal(ctx, id, params)
	assert.NoError(t, err)

	var stressLevel uint8 = 2
	params.StressLevel = &stressLevel
	err = service.PatchMeal(ctx, id, params)
	assert.NoError(t, err)

}

func TestGetFoods(t *testing.T) {
	service, ctx := initAPITest(t)

	resp, err := service.GetFoods(ctx, 1)
	assert.NoError(t, err)

	id, cleanup := service.createTestMeal(t)
	defer cleanup()

	resp, err = service.GetFoods(ctx, id)
	assert.NoError(t, err)
	assert.Len(t, resp.Foods, 0)
}

func TestPostFood(t *testing.T) {
	service, ctx := initAPITest(t)

	mealId, cleanup := service.createTestMeal(t)
	defer cleanup()

	var params = FoodParams{IngredientName: "New", Condition: entity.Cooked}
	food, err := service.PostFood(ctx, mealId, params)
	defer service.DeleteFood(ctx, food.Food.ID)
	assert.NoError(t, err)
}

func TestDeleteFoodAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	_, err := service.DeleteFood(ctx, 100)
	assert.EqualError(t, err, "not_found: not found")

	mealId, cleanup := service.createTestMeal(t)
	defer cleanup()
	var params = FoodParams{IngredientName: "New", Condition: entity.Cooked}
	food, _ := service.PostFood(ctx, mealId, params)

	ing, err := service.DeleteFood(ctx, food.Food.ID)
	assert.NoError(t, err)
	assert.Len(t, ing.Ingredients, 0)
}
