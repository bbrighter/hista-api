package meals

import (
	"context"
	_ "embed"
	"log"
	"testing"

	"encore.dev"
	"github.com/stretchr/testify/assert"
)

func initAPITest(t *testing.T) (*Service, context.Context) {
	var ctx context.Context = context.TODO()
	service := initTest(t)
	return service, ctx
}

//go:embed fixtures.sql
var fixtures string

func initTest(t *testing.T) *Service {
	service, err := initService()
	assert.NoError(t, err)
	if encore.Meta().Environment.Cloud == encore.CloudLocal {
		if _, err := histaDB.Exec(context.Background(), fixtures); err != nil {
			log.Fatalln("unable to add fixtures:", err)
		}
	}
	return service
}
