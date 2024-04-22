package meals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplaceIngredient(t *testing.T) {
	service, _ := initService()
	var ing Ingredient
	var err error

	ing, err = service.createOrReplaceIngredient("Name")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, ing.ID)
	assert.EqualValues(t, "Name", ing.Name)

	// Verify idempotency
	ing, err = service.createOrReplaceIngredient("Name")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, ing.ID)

	// Verify that new items get new names
	ing, err = service.createOrReplaceIngredient("Name2")
	assert.NoError(t, err)
	assert.EqualValues(t, 2, ing.ID)
}
