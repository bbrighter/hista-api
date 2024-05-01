package meals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToIngredientResponse(t *testing.T) {
	t.Parallel()
	var ingredient = Ingredient{
		ID:   1,
		Name: "Name",
	}
	var resp IngredientResponse = ingredient.toIngredientResponse()
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, "Name", resp.Name)
}

func TestToIngredientsResponse(t *testing.T) {
	t.Parallel()
	var ingredient1 = Ingredient{ID: 1, Name: "Name"}
	var ingredient2 = Ingredient{ID: 2, Name: "Name 2"}
	var ingredients = Ingredients{ingredient1, ingredient2}
	var resp IngredientsResponse = ingredients.toIngredientsResponse()
	assert.Len(t, resp.Ingredients, 2)
	assert.EqualValues(t, 1, resp.Ingredients[0].ID)
	assert.Equal(t, "Name", resp.Ingredients[0].Name)
}

func TestFoodsToFoodsResponse(t *testing.T) {
	t.Parallel()
	var food1 = Food{
		ID:           1,
		Ingredient:   Ingredient{ID: 10, Name: "Name"},
		IngredientID: 10,
		Condition:    Cooked,
		MealID:       100}
	var foods = Foods{food1}

	var resps []FoodResponse = foods.toFoodsResponse()

	assert.Len(t, resps, 1)
	var resp FoodResponse = resps[0]
	assert.Equal(t, Cooked, resp.Condition)
	assert.EqualValues(t, resp.ID, 1)
	assert.EqualValues(t, resp.Ingredient.ID, 10)
	assert.Equal(t, resp.Ingredient.Name, "Name")
}

func TestToMealMetaResponse(t *testing.T) {
	t.Parallel()

	var meal = Meal{
		ID:    1,
		Date:  time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC),
		Foods: []Food{{ID: 1}},
	}
	var resp MealMetaResponse = meal.toMealMetaResponse()
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC), resp.Date)
}

func TestToMealResponse(t *testing.T) {
	t.Parallel()

	var meal = Meal{
		ID:   1,
		Date: time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC),
		Foods: []Food{{
			ID:           1,
			Ingredient:   Ingredient{ID: 10, Name: "Ingredient"},
			IngredientID: 10,
			Condition:    Cooked,
			MealID:       2,
		}},
	}

	var resp MealResponse = meal.toMealResponse()

	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC), resp.Date)
	assert.Len(t, resp.Foods, 1)
	food := resp.Foods[0]
	assert.EqualValues(t, 1, food.ID)
	assert.Equal(t, Cooked, food.Condition)
	assert.EqualValues(t, 10, food.Ingredient.ID)
	assert.Equal(t, "Ingredient", food.Ingredient.Name)
}
