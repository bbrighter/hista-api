package api

import (
	"context"
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

var testFood *entity.Food = new(entity.Food)
var testIngredient *entity.Ingredient = new(entity.Ingredient)
var testMeal *entity.Meal = new(entity.Meal)

func (service *Service) createTestFood(t *testing.T) func(t *testing.T) {
	meal, err := service.meals.Create(nil)
	id := meal.ID
	testMeal.ID = id
	assert.NoError(t, err)
	food, ings, err := service.foods.Create(id, "ingredient", 0)
	testFood = &food
	testIngredient = &ings[0]
	assert.NoError(t, err)
	cleanup := func(t *testing.T) {
		ctx := context.TODO()
		_, err := service.DeleteMeal(ctx, id)
		assert.NoError(t, err)
		testFood = new(entity.Food)
		testIngredient = new(entity.Ingredient)
		testMeal = new(entity.Meal)
	}
	return cleanup
}

func TestDeleteFood(t *testing.T) {
	service, ctx := initAPITest(t)

	_, err := service.DeleteFood(ctx, 1)
	assert.EqualError(t, err, "not_found: not found")

	cleanup := service.createTestFood(t)
	defer cleanup(t)

	ing, err := service.DeleteFood(ctx, testFood.ID)
	assert.NoError(t, err)
	assert.Len(t, ing.Ingredients, 0)
}

func TestPatchFoodCondition(t *testing.T) {
	service, ctx := initAPITest(t)

	var params = FoodConditionParams{Condition: "raw"}

	err := service.PatchFoodCondition(ctx, 100, params)
	assert.EqualError(t, err, "not_found: not found")

	cleanup := service.createTestFood(t)
	defer cleanup(t)

	err = service.PatchFoodCondition(ctx, testFood.ID, params)
	assert.NoError(t, err)

	params.Condition = "invalid"
	err = service.PatchFoodCondition(ctx, testFood.ID, params)
	assert.EqualError(t, err, "invalid_argument: invalid condition")
}
