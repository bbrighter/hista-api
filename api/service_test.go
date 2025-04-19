package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	var err error
	_, err = initDb()

	assert.NoError(t, err)
}

func TestInitService(t *testing.T) {
	var err error
	_, err = initService()

	assert.NoError(t, err)
}

func initAPITest(t *testing.T) (*Service, context.Context) {
	service, err := initService()
	assert.NoError(t, err)
	return service, context.TODO()
}
