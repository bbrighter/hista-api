package internalAuth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginAPI(t *testing.T) {
	service, ctx := initAPITest(t)

	var params = LoginParams{
		UserName: "Test",
		Password: "TestPW",
	}

	var err error
	_, err = service.Login(ctx, params)
	assert.NoError(t, err)

	params.Password = "Wrong"
	_, err = service.Login(ctx, params)
	assert.Error(t, err)

}
