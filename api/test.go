package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initAPITest(t *testing.T) (*Service, context.Context) {
	service, err := initService()
	assert.NoError(t, err)
	return service, context.TODO()
}
