package api

import (
	"context"
	"testing"

	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const TEST_PIID_STR = "0c5e945e-ef6c-4934-91ff-702d94e2e7a8"

var TEST_PIID = uuid.FromStringOrNil(TEST_PIID_STR)

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
	ctx := t.Context()
	ctxWithVal := context.WithValue(ctx, "piid", TEST_PIID)
	sqlDb, err := et.NewTestDatabase(ctx, "hista_db")
	assert.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))
	assert.NoError(t, err)
	service := initServiceWithDb(db)

	return service, ctxWithVal
}
