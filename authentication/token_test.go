package authentication

import (
	"testing"
	"time"

	"encore.dev/types/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var GUID_STR = "54231a28-e8dc-46b8-aa7f-d75f804b2003"

func initTestToken(t *testing.T) *jwt.Token {
	guid, err := uuid.FromString(GUID_STR)
	assert.NoError(t, err)
	return GenerateToken("name", guid, []string{"app1"}, guid, time.Hour)
}

// func TestTokenProcess(t *testing.T) {
// 	token := initTestToken(t)

// 	tokenStr, err := tokenToSignedString(token, keys.privateKey)
// 	assert.NoError(t, err)

// 	claims, err := parseToken(tokenStr, keys.publicKey)
// 	assert.NoError(t, err)
// 	assert.Equal(t, GUID_STR, claims.Subject)
// 	assert.Equal(t, "name", claims.Name)
// }

// func TestParseToken(t *testing.T) {
// 	tokenStr := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmFtZSIsInN1YiI6IjU0MjMxYTI4LWU4ZGMtNDZiOC1hYTdmLWQ3NWY4MDRiMjAwMyIsImV4cCI6MTc1Nzc5NDI1NCwiaWF0IjoxNzU3NzA3ODU0fQ.lTa19mAsVl0aHPgPbBl4kR0Q2rRdAZje0GA7vJ0LpVL1-HhJQLxiSN1x3gdD8KxC93dMxBzADeTp5nzFE8QWC7sBxG9nTu5KWxgCbb1g6DOAT1BMgPob2CdmGd3IBOwnN4q2GmmPTHIbiex_k4E71V3bfGNZON6Y7Zfb5RpxzqUM2B7YsXNaZHCpBNyHlUNCpuKcoglUkavySmDTdBBkgbmUG-e_cY6Z41fcM3-cJ6FITznQefIabaBzzJ4vdCC_E13LPjmR37v9wiz1IT0J6_DfMsdVh75WkBQmlfWTJimt-H04yZlnnGwb2zaXlQWVI7AB7rQfM_NRj1o989Cstg"
// 	publicKey := NewSecrets().publicKey

// 	claims, err := ParseToken(tokenStr, publicKey)
// 	assert.NoError(t, err)
// 	sub, err := claims.GetSubject()
// 	assert.NoError(t, err)
// 	assert.Equal(t, "54231a28-e8dc-46b8-aa7f-d75f804b2003", sub)

// 	name := claims.Name
// 	assert.Equal(t, "name", name)
// }

// func TestValidateToken(t *testing.T) {
// 	tokenStr := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmFtZSIsInN1YiI6IjU0MjMxYTI4LWU4ZGMtNDZiOC1hYTdmLWQ3NWY4MDRiMjAwMyIsImV4cCI6MTc1Nzc5NDI1NCwiaWF0IjoxNzU3NzA3ODU0fQ.lTa19mAsVl0aHPgPbBl4kR0Q2rRdAZje0GA7vJ0LpVL1-HhJQLxiSN1x3gdD8KxC93dMxBzADeTp5nzFE8QWC7sBxG9nTu5KWxgCbb1g6DOAT1BMgPob2CdmGd3IBOwnN4q2GmmPTHIbiex_k4E71V3bfGNZON6Y7Zfb5RpxzqUM2B7YsXNaZHCpBNyHlUNCpuKcoglUkavySmDTdBBkgbmUG-e_cY6Z41fcM3-cJ6FITznQefIabaBzzJ4vdCC_E13LPjmR37v9wiz1IT0J6_DfMsdVh75WkBQmlfWTJimt-H04yZlnnGwb2zaXlQWVI7AB7rQfM_NRj1o989Cstg"
// 	publicKey := NewSecrets().publicKey

// 	_, err := ValidateToken(tokenStr, publicKey)
// 	assert.NoError(t, err)

// 	badTokenStr := "eyKhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmFtZSIsInN1YiI6IjU0MjMxYTI4LWU4ZGMtNDZiOC1hYTdmLWQ3NWY4MDRiMjAwMyIsImV4cCI6MTc1Nzc5NDI1NCwiaWF0IjoxNzU3NzA3ODU0fQ.lTa19mAsVl0aHPgPbBl4kR0Q2rRdAZje0GA7vJ0LpVL1-HhJQLxiSN1x3gdD8KxC93dMxBzADeTp5nzFE8QWC7sBxG9nTu5KWxgCbb1g6DOAT1BMgPob2CdmGd3IBOwnN4q2GmmPTHIbiex_k4E71V3bfGNZON6Y7Zfb5RpxzqUM2B7YsXNaZHCpBNyHlUNCpuKcoglUkavySmDTdBBkgbmUG-e_cY6Z41fcM3-cJ6FITznQefIabaBzzJ4vdCC_E13LPjmR37v9wiz1IT0J6_DfMsdVh75WkBQmlfWTJimt-H04yZlnnGwb2zaXlQWVI7AB7rQfM_NRj1o989Cstg"
// 	_, err = ValidateToken(badTokenStr, publicKey)
// 	assert.Error(t, err)
// }
