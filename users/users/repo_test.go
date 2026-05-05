package users

import (
	"context"
	"testing"

	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testPIID uuid.UUID

type RepoTestSuite struct {
	suite.Suite
	db  *gorm.DB
	tx  *gorm.DB
	ctx context.Context
	u   *UserRepo
}

func (s *RepoTestSuite) SetupSuite() {
	s.ctx = context.Background()

	sqlDb, err := et.NewTestDatabase(s.ctx, "users_db")
	s.Require().NoError(err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	s.Require().NoError(err)
	s.db = db

	testPIID, err = uuid.NewV4()
	s.Require().NoError(err)
}

func (s *RepoTestSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.Require().NoError(s.tx.Error)
	s.u = NewUserRepo(s.tx)
}

func (s *RepoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

func (s *RepoTestSuite) createTestUser() uuid.UUID {
	userId, err := uuid.NewV4()
	s.Require().NoError(err)
	user := User{ID: userId, Name: "user", Password: "password"}
	err = s.u.CreateUser(s.ctx, &user)
	s.Require().NoError(err)
	return userId
}

// func (s *RepoTestSuite) createTestUserProductInstance(userId uuid.UUID) {
// 	instance := UserProductInstance{
// 		ProductId:         "product",
// 		UserId:            userId,
// 		ProductInstanceId: testPIID,
// 	}
// 	err := gorm.G[UserProductInstance](s.tx).Create(s.ctx, &instance)
// 	s.Require().NoError(err)
// }

func (s *RepoTestSuite) createTestPermissions(userId uuid.UUID) {
	permission := UserAppPermission{
		UserProductInstanceId: 1,
		UserId:                userId,
		App:                   "app",
		UserProductInstance: UserProductInstance{
			ProductId:         "product",
			UserId:            userId,
			ProductInstanceId: testPIID,
		},
	}

	err := gorm.G[UserAppPermission](s.tx).Create(s.ctx, &permission)
	s.Require().NoError(err)
}

func (s *RepoTestSuite) uuidV4() uuid.UUID {
	id, err := uuid.NewV4()
	s.Require().NoError(err)
	return id
}

// func initTest(t *testing.T) (*UserRepo, context.Context) {
// 	ctx := context.Background()
// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	assert.NoError(t, err)
// 	err = db.AutoMigrate(
// 		&User{},
// 		&UserAppPermission{},
// 		&UserProductInstance{},
// 	)
// 	assert.NoError(t, err)
// 	db.Exec("PRAGMA foreign_keys = ON;")

// 	return NewUserRepo(db, 1), ctx
// }

func (s *RepoTestSuite) TestCreateUser() {
	id, err := uuid.NewV4()
	s.Require().NoError(err)
	user := User{ID: id, Name: "user", Password: "password"}

	err = s.u.CreateUser(s.ctx, &user)
	s.NoError(err)

	userInDb, _ := gorm.G[User](s.tx).Where("id = ?", user.ID).First(s.ctx)
	s.Equal("user", userInDb.Name)
	// s.NoError(bcrypt.CompareHashAndPassword([]byte(userInDb.Password), []byte("password")))
}

func (s *RepoTestSuite) TestDeleteUser() {
	id := s.createTestUser()

	err := s.u.DeleteUser(s.ctx, id)
	s.NoError(err)

	rows, _ := gorm.G[User](s.tx).Count(s.ctx, "*")
	s.EqualValues(rows, 0)
}

func (s *RepoTestSuite) TestDeleteUserNotFound() {
	id, _ := uuid.NewV4()
	err := s.u.DeleteUser(s.ctx, id)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *RepoTestSuite) TestUpdateUser() {
	id := s.createTestUser()

	var err error
	err = s.u.UpdateUser(s.ctx, id, map[string]any{"password": "new password"})
	s.NoError(err)

	err = s.u.UpdateUser(s.ctx, id, map[string]any{"name": "new name"})
	s.NoError(err)
}

func (s *RepoTestSuite) TestUpdateUserInvalidColumn() {
	id := s.createTestUser()

	err := s.u.UpdateUser(s.ctx, id, map[string]any{"invalid": "abc"})
	s.Error(err)
}

func (s *RepoTestSuite) TestFindUser() {
	id := s.createTestUser()

	user, err := s.u.FindUser(s.ctx, id)
	s.NoError(err)
	s.Equal(id, user.ID)
	s.Equal("user", user.Name)
	s.Equal("password", user.Password)
}

func (s *RepoTestSuite) TestFindPermissions() {
	id := s.createTestUser()
	s.createTestPermissions(id)

	permissions, err := s.u.FindPermissions(s.ctx, id)
	s.NoError(err)
	s.Len(permissions, 1)
	permission := permissions[0]
	s.Equal("app", permission.App)
	s.Equal(id, permission.UserId)
	s.Equal("product", permission.UserProductInstance.ProductId)
	s.Equal(id, permission.UserProductInstance.UserId)
}

func (s *RepoTestSuite) TestFindUserByName() {
	s.createTestUser()

	user, err := s.u.FindUserByName(s.ctx, "user")
	s.NoError(err)
	s.Equal("user", user.Name)
}

func (s *RepoTestSuite) TestFindUserByNameNotFound() {
	_, err := s.u.FindUserByName(s.ctx, "no name")
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *RepoTestSuite) TestListUsers() {
	users, err := s.u.ListUsers(s.ctx)
	s.NoError(err)
	s.Len(users, 0)

	s.createTestUser()

	users, err = s.u.ListUsers(s.ctx)
	s.NoError(err)
	s.Len(users, 1)
}

func (s *RepoTestSuite) TestAddUserToProductInstance() {
	userId := s.createTestUser()

	var user = User{ID: userId}
	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", user, []string{"app"})
	s.NoError(err)
}

func (s *RepoTestSuite) TestAddUserToProductInstanceOnlyOnce() {
	userId := s.createTestUser()

	var user = User{ID: userId}
	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", user, []string{"app"})
	s.NoError(err)
	err = s.u.AddUserToProductInstance(s.ctx, testPIID, "product", user, []string{"app"})
	s.ErrorIs(err, gorm.ErrDuplicatedKey)
}

func (s *RepoTestSuite) TestRemoveUserFromProductInstance() {
	userId := s.createTestUser()
	var user = User{ID: userId}
	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", user, []string{"app"})
	s.Require().NoError(err)

	err = s.u.RemoveUserFromProductInstance(s.ctx, testPIID, user)
	s.NoError(err)
}

func (s *RepoTestSuite) TestRemoveUserFromProductInstance_UserNotExists() {
	userId := s.createTestUser()
	var user = User{ID: userId}
	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", user, []string{"app"})
	s.Require().NoError(err)

	var otherUser = User{ID: s.uuidV4()}
	err = s.u.RemoveUserFromProductInstance(s.ctx, testPIID, otherUser)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *RepoTestSuite) TestRemoveUserFromProductInstance_InstanceNotExists() {
	userId := s.createTestUser()
	var user = User{ID: userId}
	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", user, []string{"app"})
	s.Require().NoError(err)

	otherPiid := s.uuidV4()
	err = s.u.RemoveUserFromProductInstance(s.ctx, otherPiid, user)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *RepoTestSuite) TestListUserForInstance() {
	userId := s.createTestUser()
	s.createTestPermissions(userId)

	users, err := s.u.ListUserForInstance(s.ctx, testPIID)
	s.NoError(err)
	s.Len(users, 1)
	s.Equal(userId, users[0].ID)
}

func (s *RepoTestSuite) TestListUserForInstance_InstanceNotFound() {
	userId := s.createTestUser()
	s.createTestPermissions(userId)

	_, err := s.u.ListUserForInstance(s.ctx, s.uuidV4())
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}
