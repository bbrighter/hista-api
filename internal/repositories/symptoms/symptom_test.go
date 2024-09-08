package symptoms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrReplace(t *testing.T) {
	repo := initTest(t)

	id, err := repo.CreateOrReplace("new name", 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)

	id, err = repo.CreateOrReplace("new name 2", 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, id)

	id, err = repo.CreateOrReplace("new name", 2)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, id)

	id, err = repo.CreateOrReplace("new name", 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
}
