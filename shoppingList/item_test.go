package shoppinglist

import (
	"context"
	"testing"

	"encore.app/shared/contextKeys"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const GUID_STR = "2012b8a8-df7f-407c-bda1-9567b5b8f06d"

var GUID = uuid.FromStringOrNil(GUID_STR)

func initTestService(t *testing.T) (*Service, context.Context) {
	ctx := t.Context()
	ctx = context.WithValue(ctx, contextKeys.Piid, GUID)

	sqlDb, err := et.NewTestDatabase(ctx, "shopping_list")
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))
	require.NoError(t, err)
	return initServiceWithDb(db), ctx
}

func TestItems(t *testing.T) {
	t.Skip()
	service, ctx := initTestService(t)

	resp, err := service.PostItemByName(ctx, ItemNameParams{Name: "name"})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, resp.ID)
	assert.EqualValues(t, 1, resp.ProductId)
}
