package meals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplaceIngredient(t *testing.T) {
	service, _ := initService()
	var id uint
	var err error

	id, err = service.createOrReplaceIngredient("Name")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)

	// Verify idempotency
	id, err = service.createOrReplaceIngredient("Name")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)

	// Verify that new items get new names
	id, err = service.createOrReplaceIngredient("Name2")
	assert.NoError(t, err)
	assert.EqualValues(t, 2, id)
}
