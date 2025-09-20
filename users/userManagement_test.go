package users

import (
	"context"
	"testing"

	"encore.dev/et"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initTestService(t *testing.T) (*Service, context.Context) {
	ctx := t.Context()
	sqlDb, err := et.NewTestDatabase(ctx, "users_db")
	assert.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))
	assert.NoError(t, err)

	service := setupService(db, 1)
	return service, ctx
}

func TestCRUDUser(t *testing.T) {
	service, ctx := initTestService(t)

	resp, _ := service.ListUsers(ctx)
	assert.Len(t, resp.Users, 0)

	idResp, err := service.CreateUser(ctx, UserParams{Name: "name", Password: "pw"})
	assert.NoError(t, err)
	_, err = service.CreateUser(ctx, UserParams{Name: "name", Password: "pw"})
	assert.Error(t, err)

	err = service.DeleteUser(ctx, idResp.UserId)
	assert.NoError(t, err)
	err = service.DeleteUser(ctx, idResp.UserId)
	assert.Error(t, err)
}
