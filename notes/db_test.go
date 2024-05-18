package notes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initAPITest(t *testing.T) (*Service, context.Context, func(t *testing.T)) {
	var ctx context.Context = context.TODO()
	service, teardown := initTest(t)
	return service, ctx, teardown
}

func initTest(t *testing.T) (*Service, func(t *testing.T)) {
	service, err := initService()
	assert.NoError(t, err)

	return service, service.teardown
}

func (service Service) teardown(t *testing.T) {
	var models = []interface{}{&Note{}}
	var err error
	for _, model := range models {
		err = service.db.Where("1=1").Delete(model).Error
		assert.NoError(t, err)
	}
}
