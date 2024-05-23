package meals

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initAPITest(t *testing.T) (*Service, context.Context) {
	var ctx context.Context = context.TODO()
	service := initTest(t)
	return service, ctx
}

func initTest(t *testing.T) *Service {
	service, err := initService()
	assert.NoError(t, err)
	return service
}
