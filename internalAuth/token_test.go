package internalAuth

import (
	"testing"
	"time"

	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsValidToken(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var userJulia *User
	var err error
	userJulia, err = getUserByName("Julia")
	assert.NoError(t, err)
	token, _ := userJulia.firstOrCreateValidToken(service)

	var validToken = &Token{Bearer: token.Bearer, UserID: userJulia.ID}

	var valid bool
	valid, err = validToken.isValid()
	assert.NoError(t, err)
	assert.True(t, valid)

	guid, _ := uuid.NewV4()
	var invalidToken = &Token{Bearer: guid, UserID: userJulia.ID}
	var invalid bool
	invalid, err = invalidToken.isValid()
	assert.NoError(t, err)
	assert.False(t, invalid)
}

func TestFirstOrCreateToken(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var userJulia *User
	var err error
	userJulia, err = getUserByName("Julia")
	var tok Token
	tok, err = userJulia.firstOrCreateValidToken(service)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, tok.Expires, time.Now())

	assert.Len(t, memorizedUsers, 1, "memorizedUsers")
	userJulia = memorizedUsers[0]
	assert.Len(t, userJulia.Tokens, 1, "Julia must have a token")

	// A still valid token is not re-created, but reused
	var sameToken Token
	sameToken, err = userJulia.firstOrCreateValidToken(service)
	assert.Equal(t, tok.Bearer, sameToken.Bearer)
	assert.Equal(t, tok.ID, sameToken.ID)
	assert.True(t, tok.Expires.Round(time.Second).Equal(sameToken.Expires.Round(time.Second)))
}

func TestCleanupTokens(t *testing.T) {
	service, teardown := initTest(t)
	defer teardown(t)

	var err error
	err = cleanupTokens(service)
	assert.NoError(t, err)

	// Check that an outdated token is removed
	var userID uint = memorizedUsers[0].ID
	guid, _ := uuid.NewV4()
	var token = Token{
		Bearer:  guid,
		Expires: time.Now().Add(-time.Hour),
		UserID:  userID,
	}
	err = service.db.Create(&token).Error
	assert.NoError(t, err)
	updateMemorizedUsers(service)
	assert.Len(t, memorizedUsers[0].Tokens, 1, "One token is created")

	err = cleanupTokens(service)
	assert.NoError(t, err)
	assert.Len(t, memorizedUsers[0].Tokens, 0, "Outdated token is removed")
}
