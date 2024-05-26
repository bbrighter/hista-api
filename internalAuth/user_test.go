package internalAuth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMemorizedUsers(t *testing.T) {
	service := initTest(t)

	updateMemorizedUsers(service)
	assert.Len(t, memorizedUsers, 2)
}

func TestInitializeUsers(t *testing.T) {
	service := initTest(t)

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
	var user = User{Name: "Test"}
	var isValid, isNotValid bool
	var err error
	isNotValid, err = user.isValidPassword("super")
	assert.NoError(t, err)
	assert.False(t, isNotValid)

	isValid, err = user.isValidPassword("TestPW")
	assert.NoError(t, err)
	assert.True(t, isValid)

}
