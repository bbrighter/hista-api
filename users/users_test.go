package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatchPassword(t *testing.T) {
	service, ctx := initTestService(t)

	idResp, err := service.CreateUser(ctx, UserParams{Name: "name", Password: "password"})
	require.NoError(t, err)

	err = service.PatchPassword(ctx, idResp.UserId, UserPasswordChangeParams{NewPassword: "new pw", OldPassword: "password"})
	assert.NoError(t, err)

	err = service.PatchPassword(ctx, idResp.UserId, UserPasswordChangeParams{NewPassword: "very new pw", OldPassword: "wrong pw"})
	assert.Error(t, err)
}
