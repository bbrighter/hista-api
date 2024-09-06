package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPollens(t *testing.T) {
	service, ctx := initAPITest(t)
	service.pollens.UseTestQuery(t)

	pollens, err := service.GetPollens(ctx)
	assert.NoError(t, err)
	assert.Len(t, pollens.Pollens, 0)
}

func TestUpdatePollen(t *testing.T) {
	service, ctx := initAPITest(t)
	service.pollens.UseTestQuery(t)

	err := service.UpdatePollen(ctx)
	assert.NoError(t, err)

	pollens, _ := service.GetPollens(ctx)
	assert.Len(t, pollens.Pollens, 1)

	err = service.UpdatePollen(ctx)
	assert.NoError(t, err)

	pollens, _ = service.GetPollens(ctx)
	assert.Len(t, pollens.Pollens, 1)
}
