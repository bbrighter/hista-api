package meals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplaceIngredient(t *testing.T) {
	t.Skip()
	service := initTest(t)
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
	service := initTest(t)

	var ingredients []Ingredient

	service.initData()
	ingredients = service.getIngredients()
	assert.GreaterOrEqual(t, len(ingredients), 1)
}

func TestDeleteIngredientIfUnused(t *testing.T) {
	service := initTest(t)
	service.initData()

	var err error
	err = deleteIngredientIfUnused(service.db, testIngredient.ID)
	assert.NoError(t, err)
	// Ingredient is used and should not be removed
	var rows int64
	var ingredient Ingredient
	rows = service.db.Find(&ingredient).RowsAffected
	assert.EqualValues(t, 1, rows)

	// Ingredient is not used and should be removed
	var unusedIngredient = Ingredient{Name: "Unused"}
	service.db.Create(&unusedIngredient)
	err = deleteIngredientIfUnused(service.db, unusedIngredient.ID)
	assert.NoError(t, err)

	rows = service.db.First(Ingredient{ID: unusedIngredient.ID}).RowsAffected
	assert.EqualValues(t, 0, rows)
}
