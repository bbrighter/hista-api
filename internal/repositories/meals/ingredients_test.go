package meals

import (
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

// func TestCreateOrReplaceIngredient(t *testing.T) {
// 	t.Skip()
// 	service := initTest(t)
// 	var ing entity.Ingredient
// 	var err error

// 	ing, err = service.createOrReplaceIngredient("Name")
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, "Name", ing.Name)

// 	// Verify idempotency
// 	var firstId = ing.ID
// 	ing, err = service.createOrReplaceIngredient("Name")
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, firstId, ing.ID)

// 	// Verify idempotency even when trimming
// 	ing, err = service.createOrReplaceIngredient(" Name ")
// 	assert.NoError(t, err)
// 	assert.EqualValues(t, firstId, ing.ID)

// 	// Verify that new items get new names
// 	ing, err = service.createOrReplaceIngredient("Name2")
// 	assert.NoError(t, err)
// 	assert.NotEqualValues(t, firstId, ing.ID)

// }

func TestGetIngredients(t *testing.T) {
	repo := initTest(t)

	var ingredients entity.Ingredients
	ingredients = repo.ListIngredients()
	assert.Len(t, ingredients, 0)

	repo.db.Create(&entity.Ingredient{Name: "ingredient"})
	ingredients = repo.ListIngredients()
	assert.Len(t, ingredients, 1)
}

// func TestDeleteIngredientIfUnused(t *testing.T) {
// 	service := initTest(t)

// 	service.db.Create(&entity.Food{ID: 1, Ingredient: entity.Ingredient{ID: 10}})

// 	var err error
// 	err = deleteIngredientIfUnused(service.db, 10)
// 	assert.NoError(t, err)
// 	// Ingredient is used and should not be removed
// 	var rows int64
// 	var ingredient entity.Ingredient
// 	rows = service.db.Find(&ingredient).RowsAffected
// 	assert.EqualValues(t, 1, rows)

// 	// Ingredient is not used and should be removed
// 	service.db.Create(&entity.Ingredient{ID: 1, Name: "Unused"})
// 	err = deleteIngredientIfUnused(service.db, 1)
// 	assert.NoError(t, err)

// 	rows = service.db.First(entity.Ingredient{ID: 1}).RowsAffected
// 	assert.EqualValues(t, 0, rows)
// }
