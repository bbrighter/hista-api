package internalAuth

import (
	"testing"
	"time"

	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsValidToken(t *testing.T) {
	service := initTest(t)

	var testUser *User
	var err error
	testUser, err = getUserByName("Test")
	assert.NoError(t, err)
	token, _ := testUser.firstOrCreateValidToken(service)

	var validToken = &Token{Bearer: token.Bearer, UserID: testUser.ID}

	var valid bool
	valid, err = validToken.isValid()
	assert.NoError(t, err)
	assert.True(t, valid)

	guid, _ := uuid.NewV4()
	var invalidToken = &Token{Bearer: guid, UserID: testUser.ID}
	var invalid bool
	invalid, err = invalidToken.isValid()
	assert.NoError(t, err)
	assert.False(t, invalid)
}

func TestFirstOrCreateToken(t *testing.T) {
	service := initTest(t)

	var testUser = &User{
		ID: 2,
	}
	var err error
	var tok Token
	for _, user := range memorizedUsers {
		tok, err = user.firstOrCreateValidToken(service)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, tok.Expires, time.Now())
	}

	assert.Len(t, memorizedUsers, 2, "memorizedUsers")
	for _, user := range memorizedUsers {
		assert.Len(t, user.Tokens, 1, user.Name+" must have a token")
	}

	// A still valid token is not re-created, but reused
	var sameToken Token
	sameToken, err = testUser.firstOrCreateValidToken(service)
	assert.Equal(t, tok.Bearer, sameToken.Bearer)
	assert.Equal(t, tok.ID, sameToken.ID)
	assert.True(t, tok.Expires.Round(time.Second).Equal(sameToken.Expires.Round(time.Second)))
}

func TestCleanupTokens(t *testing.T) {
	service := initTest(t)

	var err error
	err = cleanupTokens(service)
	assert.NoError(t, err)

	// Check that an outdated token is removed
	guid, _ := uuid.NewV4()
	var token = Token{
		Bearer:  guid,
		Expires: time.Now().Add(-time.Hour),
		UserID:  2, // Test user
		ID:      3,
	}
	err = service.db.Save(&token).Error
	assert.NoError(t, err)
	updateMemorizedUsers(service)
	var testUser *User
	testUser, _ = getUserByName("Test")
	assert.Len(t, testUser.Tokens, 2, "One token is created")

	err = cleanupTokens(service)
	assert.NoError(t, err)

	testUser, _ = getUserByName("Test") // need to update testuser again
	assert.Len(t, testUser.Tokens, 1, "Outdated token is removed")
}
