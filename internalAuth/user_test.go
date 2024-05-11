package internalAuth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMemorizedUsers(t *testing.T) {
	service, teardown := initTest(t)
	// service.testCreate(t)
	defer teardown(t)

	updateMemorizedUsers(service)
	assert.Len(t, memorizedUsers, 1)
}

func TestInitializeUsers(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var err error
	err = service.initializeUsers()
	assert.NoError(t, err)
}

func TestHashing(t *testing.T) {
	t.Parallel()
	result := hashing("Das ist ein Test")
	assert.Equal(t, "c423b0dcb3697b21d846dbf01f6f2438fbc973078bd8fc671551f5b98e8af0f0", result)
}

func TestIsValidPassword(t *testing.T) {
	_, teardown := initTest(t)
	defer teardown(t)

	var user = User{Name: "Julia"}
	var isValid, isNotValid bool
	var err error
	isNotValid, err = user.isValidPassword("super")
	assert.NoError(t, err)
	assert.False(t, isNotValid)

	isValid, err = user.isValidPassword("123pi")
	assert.NoError(t, err)
	assert.True(t, isValid)

}
