package users

import (
	"testing"

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
		expectError bool
	}{
		"ok":       {userName: "Name", password: "Password"},
		"no user":  {userName: "unknown", password: "Password", expectError: true},
		"wrong pw": {userName: "Name", password: "Wrong", expectError: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := service.Login(ctx, LoginParams{UserName: test.userName, Password: test.password})

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, resp.Token, "ey")
			}
		})
	}

}
