package repository

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"encore.app/entity"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var GUID_STR = "54231a28-e8dc-46b8-aa7f-d75f804b2003"

func newTestTokenRepo(t *testing.T) TokenRepo {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return TokenRepo{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

func TestGenerateToken(t *testing.T) {
	r := newTestTokenRepo(t)

	guid := uuid.FromStringOrNil(GUID_STR)
	token := r.GenerateToken("name", guid, []string{"app1"}, guid, time.Hour)

	sub, err := token.Claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, GUID_STR, sub)
	exp, err := token.Claims.GetExpirationTime()
	assert.NoError(t, err)
	assert.Less(t, exp.Time, time.Now().Add(time.Hour))
	assert.Greater(t, exp.Time, time.Now().Add(time.Minute))
	custom, ok := token.Claims.(entity.CustomClaims)
	require.True(t, ok)
	assert.Equal(t, []string{"app1"}, custom.AppIds)
	assert.Equal(t, "name", custom.Name)
	assert.Equal(t, GUID_STR, custom.PIID.String())
}
