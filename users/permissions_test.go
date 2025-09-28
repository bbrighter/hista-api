package users

import (
	"testing"

	"encore.dev/beta/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	service, ctx := initTestService(t)
	_, err := service.CreateUser(ctx, UserParams{Name: "Name", Password: "Password"})
	require.NoError(t, err)

	tests := map[string]struct {
		userName    string
		password    string
		expectError errs.ErrCode
	}{
		"ok":       {userName: "Name", password: "Password"},
		"no user":  {userName: "unknown", password: "Password", expectError: errs.NotFound},
		"wrong pw": {userName: "Name", password: "Wrong", expectError: errs.Unauthenticated},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := service.GetPermissions(ctx, LoginParams{UserName: test.userName, Password: test.password})

			if test.expectError > 0 {
				assert.Error(t, err)
				e, ok := err.(*errs.Error)
				require.True(t, ok)
				assert.Equal(t, test.expectError, e.Code)
				// assert.ErrorIs(t, err, test.expectError)

			} else {
				assert.NoError(t, err)
				// assert.Contains(t, resp.Token, "ey")
			}
		})
	}

}
