package authentication

import (
	"context"
	"strings"
	"testing"

	"encore.dev/types/uuid"
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
			var user_guid string = "00f5fd47-140c-486f-ae5a-f2926d701f30"
			var piid_guid string = "31d621bc-cba4-479e-94b6-d66f919ea612"
			ctx := context.Background()
			s, _ := initService()
			token, err := s.g.GenerateToken("name", uuid.FromStringOrNil(user_guid), []string{"app1"}, uuid.FromStringOrNil(piid_guid))
			// token := initTestToken(t)
			require.NoError(t, err)
			require.NoError(t, err)
			if !test.useOriginalToken {
				parts := strings.Split(token, ".")
				require.Len(t, parts, 3)
				tamperedPayload := parts[1][:len(parts[1])-1] + "X"
				token = parts[0] + "." + tamperedPayload + "." + parts[2]
			}
			if test.useExpiredToken {
				token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmFtZSIsInN1YiI6IjU0MjMxYTI4LWU4ZGMtNDZiOC1hYTdmLWQ3NWY4MDRiMjAwMyIsImV4cCI6MTc1Nzc5NDI1NCwiaWF0IjoxNzU3NzA3ODU0fQ.lTa19mAsVl0aHPgPbBl4kR0Q2rRdAZje0GA7vJ0LpVL1-HhJQLxiSN1x3gdD8KxC93dMxBzADeTp5nzFE8QWC7sBxG9nTu5KWxgCbb1g6DOAT1BMgPob2CdmGd3IBOwnN4q2GmmPTHIbiex_k4E71V3bfGNZON6Y7Zfb5RpxzqUM2B7YsXNaZHCpBNyHlUNCpuKcoglUkavySmDTdBBkgbmUG-e_cY6Z41fcM3-cJ6FITznQefIabaBzzJ4vdCC_E13LPjmR37v9wiz1IT0J6_DfMsdVh75WkBQmlfWTJimt-H04yZlnnGwb2zaXlQWVI7AB7rQfM_NRj1o989Cstg"
			}

			id, _, err := AuthHandler(ctx, &AuthParams{Token: token})
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, user_guid, id)
			}
		})
	}
}
