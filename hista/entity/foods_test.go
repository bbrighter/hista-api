package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoodToFoodResponse(t *testing.T) {
	t.Parallel()
	var amount int = 50
	var food = Food{
		ID:           1,
		Ingredient:   Ingredient{ID: 10, Name: "Name"},
		IngredientID: 10,
		Condition:    Cooked,
		MealID:       100,
		Amount:       &amount,
	}
	var resp FoodResponse = food.ToFoodResponse()

	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, uint(10), resp.Ingredient.ID)
	assert.Equal(t, "Name", resp.Ingredient.Name)
	assert.Equal(t, Cooked, resp.Condition)
	assert.Equal(t, 50, *resp.Amount)
}

func TestFoodsToFoodsResponse(t *testing.T) {
	t.Parallel()
	var amount int = 50
	var food1 = &Food{
		ID:           1,
		Ingredient:   Ingredient{ID: 10, Name: "Name"},
		IngredientID: 10,
		Condition:    Cooked,
		MealID:       100,
		Amount:       &amount,
	}
	var foods = Foods{food1}

	var resps []FoodResponse = foods.ToFoodsResponse()

	assert.Len(t, resps, 1)
	var resp FoodResponse = resps[0]
	assert.Equal(t, Cooked, resp.Condition)
	assert.EqualValues(t, resp.ID, 1)
	assert.EqualValues(t, resp.Ingredient.ID, 10)
	assert.Equal(t, resp.Ingredient.Name, "Name")
	assert.Equal(t, 50, *resp.Amount)
}
