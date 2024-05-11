package internalAuth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginAPI(t *testing.T) {
	service, ctx, teardown := initAPITest(t)
	service.useTestPassword(t)
	defer teardown(t)

	var user User
	service.db.First(&user)

	var params = LoginParams{
		UserName: "Julia",
		Password: "TestPW",
	}

	var err error
	_, err = service.Login(ctx, params)
	assert.NoError(t, err)

	params.Password = "Wrong"
	_, err = service.Login(ctx, params)
	assert.Error(t, err)

}
