package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (service *Service) createTestFood(t *testing.T) (uint, func(t *testing.T)) {
	id, err := service.mealUC.CreateMeal(nil)
	assert.NoError(t, err)
	food, _, err := service.food.CreateFood(id, "ingredient", 0)
	assert.NoError(t, err)
	cleanup := func(t *testing.T) {
		ctx := context.TODO()
		err := service.DeleteMeal(ctx, id)
		assert.NoError(t, err)
	}
	return food.ID, cleanup
}

func TestDeleteFood(t *testing.T) {
	service, ctx := initAPITest(t)

	_, err := service.DeleteFood(ctx, 1)
	assert.EqualError(t, err, "not_found: not found")

	id, cleanup := service.createTestFood(t)
	defer cleanup(t)

	ing, err := service.DeleteFood(ctx, id)
	assert.NoError(t, err)
	assert.Len(t, ing.Ingredients, 0)
}

func TestPatchFoodCondition(t *testing.T) {
	service, ctx := initAPITest(t)

	var params = FoodConditionParams{Condition: "raw"}

	err := service.PatchFoodCondition(ctx, 100, params)
	assert.EqualError(t, err, "not_found: not found")

	id, cleanup := service.createTestFood(t)
	defer cleanup(t)

	err = service.PatchFoodCondition(ctx, id, params)
	assert.NoError(t, err)

	params.Condition = "invalid"
	err = service.PatchFoodCondition(ctx, id, params)
	assert.EqualError(t, err, "invalid_argument: invalid condition")
}
