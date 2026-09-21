package userproductinstance

import (
	"context"
	"testing"

	"encore.app/users/internal/shared"
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
	u   *userProductInstanceRepo
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
	s.u = NewUserProductInstanceRepo(s.tx)
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
	user := shared.User{ID: userId, Name: "user", Password: "password"}
	err = gorm.G[shared.User](s.tx).Create(s.ctx, &user)
	s.Require().NoError(err)
	return userId
}

func (s *RepoTestSuite) createTestPermissions(userId uuid.UUID) {
	permission := shared.UserAppPermission{
		UserProductInstanceId: 1,
		UserId:                userId,
		App:                   "app",
		UserProductInstance: shared.UserProductInstance{
			ProductId:         "product",
			UserId:            userId,
			ProductInstanceId: testPIID,
		},
	}

	err := gorm.G[shared.UserAppPermission](s.tx).Create(s.ctx, &permission)
	s.Require().NoError(err)
}

func (s *RepoTestSuite) uuidV4() uuid.UUID {
	id, err := uuid.NewV4()
	s.Require().NoError(err)
	return id
}

func (s *RepoTestSuite) TestAddUserToProductInstance() {
	userId := s.createTestUser()

	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", userId, []string{"app"})
	s.NoError(err)
}

func (s *RepoTestSuite) TestAddUserToProductInstanceOnlyOnce() {
	userId := s.createTestUser()

	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", userId, []string{"app"})
	s.NoError(err)
	err = s.u.AddUserToProductInstance(s.ctx, testPIID, "product", userId, []string{"app"})
	s.ErrorIs(err, gorm.ErrDuplicatedKey)
}

func (s *RepoTestSuite) TestRemoveUserFromProductInstance() {
	userId := s.createTestUser()

	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", userId, []string{"app"})
	s.Require().NoError(err)

	err = s.u.RemoveUserFromProductInstance(s.ctx, testPIID, userId)
	s.NoError(err)
}

func (s *RepoTestSuite) TestRemoveUserFromProductInstance_UserNotExists() {
	userId := s.createTestUser()

	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", userId, []string{"app"})
	s.Require().NoError(err)

	err = s.u.RemoveUserFromProductInstance(s.ctx, testPIID, s.uuidV4())
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *RepoTestSuite) TestRemoveUserFromProductInstance_InstanceNotExists() {
	userId := s.createTestUser()

	err := s.u.AddUserToProductInstance(s.ctx, testPIID, "product", userId, []string{"app"})
	s.Require().NoError(err)

	otherPiid := s.uuidV4()
	err = s.u.RemoveUserFromProductInstance(s.ctx, otherPiid, userId)
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
