package authentication

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler(t *testing.T) {
	tests := map[string]struct {
		useOriginalToken bool
		useExpiredToken  bool
		expectError      bool
	}{
		"ok":             {useOriginalToken: true},
		"tampered token": {useOriginalToken: false, expectError: true},
		"expired token":  {useExpiredToken: true, expectError: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			token := initTestToken(t)
			privateKey := keys.privateKey
			signedToken, err := token.SignedString(privateKey)
			require.NoError(t, err)
			if !test.useOriginalToken {
				parts := strings.Split(signedToken, ".")
				require.Len(t, parts, 3)
				tamperedPayload := parts[1][:len(parts[1])-1] + "X"
				signedToken = parts[0] + "." + tamperedPayload + "." + parts[2]
			}
			if test.useExpiredToken {
				signedToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmFtZSIsInN1YiI6IjU0MjMxYTI4LWU4ZGMtNDZiOC1hYTdmLWQ3NWY4MDRiMjAwMyIsImV4cCI6MTc1Nzc5NDI1NCwiaWF0IjoxNzU3NzA3ODU0fQ.lTa19mAsVl0aHPgPbBl4kR0Q2rRdAZje0GA7vJ0LpVL1-HhJQLxiSN1x3gdD8KxC93dMxBzADeTp5nzFE8QWC7sBxG9nTu5KWxgCbb1g6DOAT1BMgPob2CdmGd3IBOwnN4q2GmmPTHIbiex_k4E71V3bfGNZON6Y7Zfb5RpxzqUM2B7YsXNaZHCpBNyHlUNCpuKcoglUkavySmDTdBBkgbmUG-e_cY6Z41fcM3-cJ6FITznQefIabaBzzJ4vdCC_E13LPjmR37v9wiz1IT0J6_DfMsdVh75WkBQmlfWTJimt-H04yZlnnGwb2zaXlQWVI7AB7rQfM_NRj1o989Cstg"
			}

			id, _, err := AuthHandler(ctx, &AuthParams{token: signedToken})
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, "54231a28-e8dc-46b8-aa7f-d75f804b2003", id)
			}
		})
	}
}
