package meals

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFoods(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	context := context.TODO()

	resp, err := service.GetFoods(context, 1)
	assert.NoError(t, err)
	assert.Len(t, resp.Foods, 0)
}

func TestPostFood(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	ctx := context.TODO()

	var params = FoodParams{IngredientName: "New", Condition: Cooked}
	_, err := service.PostFood(ctx, 1, params)
	assert.Error(t, err)
}

func TestDeleteFoodAPI(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	ctx := context.TODO()

	var err error
	err = service.DeleteFood(ctx, 1, 1)
	assert.Error(t, err)

	service.testCreateMeal(t)

	err = service.DeleteFood(ctx, 1, 1)
	assert.NoError(t, err)
}

func TestPatchFoodCondition(t *testing.T) {
	service, teardown := initTest(t)
	ctx := context.TODO()
	service.testCreateMeal(t)
	defer teardown(t)

	var params = FoodConditionParams{Condition: "raw"}

	var err error
	service.PatchFoodCondition(ctx, 1, 1, params)
	assert.NoError(t, err)

	var invalidParams = FoodConditionParams{Condition: "not valid"}
	service.PatchFoodCondition(ctx, 1, 1, invalidParams)
	assert.NoError(t, err)
}
