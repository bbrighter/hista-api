package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToIngredientResponse(t *testing.T) {
	t.Parallel()
	var prot int = 100
	var carb int = 20
	var fat int = 15
	var fiber int = 0
	var ingredient = Ingredient{
		ID:         1,
		Name:       "Name",
		IsArchived: true,
		Nutrition: Nutrition{
			Protein:      &prot,
			Carbohydrate: &carb,
			Fat:          &fat,
			Fiber:        &fiber,
		},
	}
	var resp IngredientResponse = ingredient.ToIngredientResponse()
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, "Name", resp.Name)
	assert.Equal(t, true, resp.IsArchived)
	assert.Equal(t, prot, resp.Nutrition.Protein)
	assert.Equal(t, carb, resp.Nutrition.Carbohydrate)
	assert.Equal(t, fat, resp.Nutrition.Fat)
	assert.Equal(t, fiber, resp.Nutrition.Fiber)

}

func TestToIngredientsResponse(t *testing.T) {
	t.Parallel()
	var ingredient1 = &Ingredient{ID: 1, Name: "Name"}
	var ingredient2 = &Ingredient{ID: 2, Name: "Name 2"}
	var ingredients = Ingredients{ingredient1, ingredient2}
	var resp IngredientsResponse = ingredients.ToIngredientsResponse()
	assert.Len(t, resp.Ingredients, 2)
	assert.EqualValues(t, 1, resp.Ingredients[0].ID)
	assert.Equal(t, "Name", resp.Ingredients[0].Name)
}
