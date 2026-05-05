package users

import (
	"context"
	"testing"

	"encore.dev/beta/errs"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type apiTestSuite struct {
	suite.Suite
	ctx     context.Context
	service *Service
}

func (s *apiTestSuite) SetupTest() {
	s.ctx = context.Background()
	sqlDb, err := et.NewTestDatabase(s.ctx, "users_db")
	s.Require().NoError(err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))
	s.Require().NoError(err)
	cost := 4
	s.service = setupService(db, cost)
}

func (s *apiTestSuite) assertErrCode(err error, code errs.ErrCode) {
	e, ok := err.(*errs.Error)
	s.True(ok)
	s.Equal(code, e.Code)
}

func (s *apiTestSuite) createTestUser() uuid.UUID {
	idResp, err := s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "pw"})
	s.Require().NoError(err)
	return idResp.UserId
}

func (s *apiTestSuite) uuid() uuid.UUID {
	id, err := uuid.NewV4()
	s.Require().NoError(err)
	return id
}

func TestApi(t *testing.T) {
	suite.Run(t, new(apiTestSuite))
}

func (s *apiTestSuite) TestCreateUser() {
	idResp, err := s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "pw"})
	s.NoError(err)

	// Cannot create twice
	_, err = s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "pw"})
	s.Error(err)
	s.assertErrCode(err, errs.AlreadyExists)
	existsResp, err := s.service.Exists(s.ctx, "name")
	s.NoError(err)
	s.Equal(idResp, existsResp)
}

func (s *apiTestSuite) TestListUsers() {
	id := s.createTestUser()

	resp, err := s.service.ListUsers(s.ctx)
	s.NoError(err)
	s.Len(resp.Users, 1)
	s.Equal(UserResponse{ID: id, Name: "name"}, resp.Users[0])
}

func (s *apiTestSuite) TestDeleteUser() {
	id := s.createTestUser()

	err := s.service.DeleteUser(s.ctx, id)
	s.NoError(err)

	_, err = s.service.Exists(s.ctx, "name")
	s.Error(err)
}

func (s *apiTestSuite) TestDeleteUserNotFound() {
	id := s.uuid()
	err := s.service.DeleteUser(s.ctx, id)
	s.Error(err)
	s.assertErrCode(err, errs.NotFound)
}

func (s *apiTestSuite) TestManageUserToProductInstance() {
	userId := s.createTestUser()
	piid := s.uuid()
	apps := []string{"app"}
	productId := "product"

	// Add user to product instance
	err := s.service.AddUserToProductInstance(s.ctx, userId, piid, AddUserToProductInstanceParams{AppIds: apps, ProductId: productId})

	s.NoError(err)
	resp, err := s.service.ListUsersForProductInstance(s.ctx, piid)
	s.NoError(err)
	s.Len(resp.Users, 1)
	user := resp.Users[0]
	s.Equal(UserResponse{ID: userId, Name: "name"}, user)

	// Remove again
	err = s.service.RemoveUserFromProductInstance(s.ctx, userId, piid)

	s.NoError(err)
	_, err = s.service.ListUsersForProductInstance(s.ctx, piid)
	s.Error(err)
	s.assertErrCode(err, errs.NotFound)
}

func (s *apiTestSuite) TestGetPermissions() {
	_, err := s.service.CreateUser(s.ctx, UserParams{Name: "Name", Password: "Password"})
	s.Require().NoError(err)

	tests := map[string]struct {
		userName    string
		password    string
		expectError errs.ErrCode
	}{
		"ok":       {userName: "Name", password: "Password"},
		"no user":  {userName: "unknown", password: "Password", expectError: errs.NotFound},
		"wrong pw": {userName: "Name", password: "Wrong", expectError: errs.Unauthenticated},
	}

	for name, test := range tests {
		s.Run(name, func() {
			_, err := s.service.GetPermissions(s.ctx, LoginParams{UserName: test.userName, Password: test.password})

			if test.expectError > 0 {
				s.assertErrCode(err, test.expectError)
				return
			}
			s.NoError(err)
		})
	}
}

func (s *apiTestSuite) TestPatchPassword() {
	idResp, err := s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "password"})
	s.NoError(err)

	err = s.service.PatchPassword(s.ctx, idResp.UserId, UserPasswordChangeParams{NewPassword: "new pw", OldPassword: "password"})
	s.NoError(err)

	err = s.service.PatchPassword(s.ctx, idResp.UserId, UserPasswordChangeParams{NewPassword: "very new pw", OldPassword: "wrong pw"})
	s.Error(err)
	s.assertErrCode(err, errs.Unauthenticated)
}
