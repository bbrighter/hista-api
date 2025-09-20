package repository

import (
	"context"
	"testing"

	"encore.app/users/entity"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initTest(t *testing.T) (*UserRepo, context.Context) {
	ctx := context.Background()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(
		&entity.User{},
		&entity.UserAppPermission{},
		&entity.UserProductInstance{},
	)
	assert.NoError(t, err)
	db.Exec("PRAGMA foreign_keys = ON;")

	return NewUserRepo(db, 1), ctx
}

func TestCreateUser(t *testing.T) {
	r, ctx := initTest(t)

	user, err := r.Create(ctx, "user", "password")
	assert.NoError(t, err)

	userInDb, _ := gorm.G[entity.User](r.db).Where("id = ?", user.ID).First(ctx)

	assert.Equal(t, "user", userInDb.Name)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(userInDb.Password), []byte("password")))
}

func TestDeleteUser(t *testing.T) {
	tests := map[string]struct {
		useUserId             bool
		expectedError         bool
		expectedNumberOfUsers int64
	}{
		"ok":        {useUserId: true, expectedNumberOfUsers: 0},
		"not found": {useUserId: false, expectedError: true, expectedNumberOfUsers: 1},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, ctx := initTest(t)
			user, err := r.Create(ctx, "user", "password")
			assert.NoError(t, err)

			var userId uuid.UUID
			if test.useUserId {
				userId = user.ID
			} else {
				userId, err = uuid.NewV4()
				assert.NoError(t, err)
			}

			err = r.Delete(ctx, userId)

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	tests := map[string]struct {
		useUserId     bool
		expectedError bool
	}{
		"ok":        {useUserId: true},
		"not found": {useUserId: false, expectedError: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, ctx := initTest(t)
			user, err := r.Create(ctx, "user", "password")
			assert.NoError(t, err)

			var userId uuid.UUID
			if test.useUserId {
				userId = user.ID
			} else {
				userId, err = uuid.NewV4()
				assert.NoError(t, err)
			}

			err = r.ChangePassword(ctx, userId, "new password")

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				newUser, _ := gorm.G[entity.User](r.db).Where("id = ?", user.ID).First(ctx)
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte("new password")))
			}
		})
	}
}

func TestLogin(t *testing.T) {

	tests := map[string]struct {
		userName      string
		password      string
		expectedError bool
	}{
		"ok":             {userName: "name", password: "password"},
		"user not found": {userName: "unknown", password: "password", expectedError: true},
		"wrong password": {userName: "name", password: "wrong password", expectedError: true},
	}

	r, ctx := initTest(t)
	user, err := r.Create(ctx, "name", "password")
	require.NoError(t, err)
	guid, err := uuid.NewV4()
	require.NoError(t, err)
	err = r.AddUser(ctx, guid, "product-id", user, []string{"app1"})
	require.NoError(t, err)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			user, perm, err := r.Login(ctx, test.userName, test.password)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.userName, user.Name)
				assert.Len(t, perm, 1)
				assert.Equal(t, guid, perm[0].UserProductInstance.ProductInstanceId)
				assert.EqualValues(t, "app1", perm[0].App)
			}
		})
	}
}

func TestFindUser(t *testing.T) {

	tests := map[string]struct {
		useId         bool
		expectedError bool
	}{
		"ok":             {useId: true},
		"user not found": {useId: false, expectedError: true},
	}

	r, ctx := initTest(t)
	testUser, err := r.Create(ctx, "name", "password")
	assert.NoError(t, err)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var guid uuid.UUID
			if test.useId {
				guid = testUser.ID
			} else {
				guid, _ = uuid.NewV4()
			}

			user, err := r.Find(ctx, guid)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, guid, user.ID)
				assert.Equal(t, "name", user.Name)
			}
		})
	}
}

func TestAddUser(t *testing.T) {
	r, ctx := initTest(t)

	userId, err := uuid.NewV4()
	assert.NoError(t, err)

	var user = entity.User{ID: userId, Name: "name", Password: "pw"}
	err = gorm.G[entity.User](r.db).Create(ctx, &user)
	assert.NoError(t, err)

	tests := map[string]struct {
		useExistingUser bool
		expectedError   bool
	}{
		"ok":             {useExistingUser: true},
		"user not found": {useExistingUser: false, expectedError: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, err)
			if !test.useExistingUser {
				otherGuid, err := uuid.NewV4()
				assert.NoError(t, err)
				user = entity.User{ID: otherGuid}
			}

			err = r.AddUser(ctx, uuid.UUID{}, "prod", user, []string{"app1", "app2"})

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveUser(t *testing.T) {
	invalidUserId, _ := uuid.NewV4()
	piid, _ := uuid.NewV4()
	invalidPiid, _ := uuid.NewV4()

	tests := map[string]struct {
		useExistingUser     bool
		useExistingInstance bool
		expectError         bool
	}{
		"ok":                   {useExistingUser: true, useExistingInstance: true},
		"user not found":       {useExistingUser: false, useExistingInstance: true, expectError: true},
		"user not in instance": {useExistingUser: true, useExistingInstance: false, expectError: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, ctx := initTest(t)
			user, err := r.Create(ctx, "name", "password")
			assert.NoError(t, err)
			err = r.AddUser(ctx, piid, "prod", user, []string{"app1", "app2"})
			assert.NoError(t, err)

			var usePiid uuid.UUID = piid
			if !test.useExistingInstance {
				usePiid = invalidPiid
			}
			var useUserId uuid.UUID = user.ID
			if !test.useExistingUser {
				useUserId = invalidUserId
			}
			err = r.RemoveUser(ctx, usePiid, entity.User{ID: useUserId})

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
