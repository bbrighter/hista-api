package api

import (
	"context"
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (service *Service) createTestFood(ctx context.Context, t *testing.T) (foodId uint, ingredientId uint) {
	meal, err := service.PostMeal(ctx, entity.MealParams{})
	require.NoError(t, err)
	foodResp, err := service.PostFood(ctx, meal.ID, FoodParams{IngredientName: "ingredient", IngredientID: 0})
	require.NoError(t, err)
	return foodResp.Food.ID, foodResp.Food.Ingredient.ID
}

func TestDeleteFood(t *testing.T) {
	service, ctx := initAPITest(t)

	_, err := service.DeleteFood(ctx, 1)
	assert.EqualError(t, err, "not_found: not found")

	foodId, _ := service.createTestFood(ctx, t)

	ing, err := service.DeleteFood(ctx, foodId)
	assert.NoError(t, err)
	assert.Len(t, ing.Ingredients, 0)
}

func TestPatchFoodCondition(t *testing.T) {
	service, ctx := initAPITest(t)

	var params = FoodConditionParams{Condition: "raw"}

	err := service.PatchFoodCondition(ctx, 100, params)
	assert.EqualError(t, err, "not_found: not found")

	foodId, _ := service.createTestFood(ctx, t)

	err = service.PatchFoodCondition(ctx, foodId, params)
	assert.NoError(t, err)

	params.Condition = "invalid"
	err = service.PatchFoodCondition(ctx, foodId, params)
	assert.EqualError(t, err, "invalid_argument: invalid condition")
}
