package meals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplaceIngredient(t *testing.T) {
	t.Skip()
	service, teardown := initTest(t)
	defer teardown(t)
	var ing Ingredient
	var err error

	ing, err = service.createOrReplaceIngredient("Name")
	assert.NoError(t, err)
	assert.EqualValues(t, "Name", ing.Name)

	// Verify idempotency
	var firstId = ing.ID
	ing, err = service.createOrReplaceIngredient("Name")
	assert.NoError(t, err)
	assert.EqualValues(t, firstId, ing.ID)

	// Verify that new items get new names
	ing, err = service.createOrReplaceIngredient("Name2")
	assert.NoError(t, err)
	assert.NotEqualValues(t, firstId, ing.ID)
}

func TestGetIngredients(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var ingredients []Ingredient

	ingredients = service.getIngredients()
	assert.Len(t, ingredients, 0)

	service.createOrReplaceIngredient("New")
	ingredients = service.getIngredients()
	assert.Len(t, ingredients, 1)
}

func TestDeleteIngredientIfUnused(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var meal Meal = service.testCreateMeal(t)
	var ingredient Ingredient = meal.Foods[0].Ingredient

	var err error
	err = deleteIngredientIfUnused(service.db, ingredient.ID)
	assert.NoError(t, err)
	// Ingredient is used and should not be removed
	var rows int64
	rows = service.db.Find(&ingredient).RowsAffected
	assert.EqualValues(t, 1, rows)

	// Ingredient is not used and should be removed
	service.db.Create(&Ingredient{ID: 2, Name: "Name 2"})
	err = deleteIngredientIfUnused(service.db, 2)
	assert.NoError(t, err)

	rows = service.db.Find(&Ingredient{ID: 2}).RowsAffected
	assert.EqualValues(t, 0, rows)
}
