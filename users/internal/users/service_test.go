package users

import (
	"context"
	"errors"
	"testing"

	"encore.app/users/internal/shared"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockedUserRepo struct {
	mock.Mock
}

func (m *mockedUserRepo) CreateUser(ctx context.Context, user *shared.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *mockedUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockedUserRepo) UpdateUser(ctx context.Context, id uuid.UUID, values map[string]any) error {
	args := m.Called(ctx, id, values)
	return args.Error(0)
}
func (m *mockedUserRepo) FindUser(ctx context.Context, id uuid.UUID) (shared.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(shared.User), args.Error(1)
}
func (m *mockedUserRepo) FindPermissions(ctx context.Context, id uuid.UUID) ([]shared.UserAppPermission, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]shared.UserAppPermission), args.Error(1)
}
func (m *mockedUserRepo) FindUserByName(ctx context.Context, name string) (shared.User, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(shared.User), args.Error(1)
}
func (m *mockedUserRepo) ListUsers(ctx context.Context) ([]shared.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]shared.User), args.Error(1)
}

// func (m *mockedUserRepo) AddUserToProductInstance(ctx context.Context, instanceId uuid.UUID, productId string, user User, apps []string) error {
// 	args := m.Called(ctx, instanceId, productId, user, apps)
// 	return args.Error(0)
// }
// func (m *mockedUserRepo) RemoveUserFromProductInstance(ctx context.Context, instanceId uuid.UUID, user User) error {
// 	args := m.Called(ctx, instanceId, user)
// 	return args.Error(0)
// }
// func (m *mockedUserRepo) ListUserForInstance(ctx context.Context, instanceId uuid.UUID) ([]User, error) {
// 	args := m.Called(ctx, instanceId)
// 	return args.Get(0).([]User), args.Error(1)
// }

type mockedEncryption struct {
	mock.Mock
}

func (m *mockedEncryption) GeneratePassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
func (m *mockedEncryption) ValidatePassword(password1, password2 string) error {
	args := m.Called(password1, password2)
	return args.Error(0)
}

type serviceTestSuite struct {
	suite.Suite
	service *UserService
	ctx     context.Context
	r       *mockedUserRepo
	e       *mockedEncryption
}

func (s *serviceTestSuite) SetupTest() {
	s.r = new(mockedUserRepo)
	s.e = new(mockedEncryption)

	s.ctx = s.T().Context()

	s.service = &UserService{u: s.r, e: s.e}
}

func (s *serviceTestSuite) SetupSubTest() {
	s.r = new(mockedUserRepo)
	s.e = new(mockedEncryption)

	s.ctx = s.T().Context()

	s.service = &UserService{u: s.r, e: s.e}
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (s *serviceTestSuite) TestCreateUser() {
	s.e.On("GeneratePassword", "password").Return("xyz", nil)
	s.r.On("CreateUser", s.ctx, mock.MatchedBy(func(u *shared.User) bool {
		return u.Name == "name" && u.Password == "xyz"
	})).Return(nil)

	user, err := s.service.CreateUser(s.ctx, "name", "password")
	s.NoError(err)
	s.NotNil(user)
	s.Equal("name", user.Name)
	s.r.AssertExpectations(s.T())
}

func (s *serviceTestSuite) TestCreateUser_PasswordFailed() {
	s.e.On("GeneratePassword", "password").Return("", errors.New("pw generation failed"))

	_, err := s.service.CreateUser(s.ctx, "name", "password")
	s.Error(err)
	s.r.AssertExpectations(s.T())
}

func (s *serviceTestSuite) TestChangePassword() {
	userId, err := uuid.NewV4()
	s.Require().NoError(err)

	var someErr = errors.New("some err")
	tests := map[string]struct {
		user          shared.User
		findErr       error
		validationErr error
		updateErr     error
		expectedError error
	}{
		"ok":               {user: shared.User{Password: "xxx"}},
		"user not found":   {findErr: someErr, expectedError: someErr},
		"password invalid": {user: shared.User{Password: "xxx"}, validationErr: someErr, expectedError: someErr},
		"update failed":    {user: shared.User{Password: "xxx"}, updateErr: someErr, expectedError: someErr},
	}
	for name, test := range tests {
		s.Run(name, func() {
			s.r.On("FindUser", s.ctx, userId).Return(test.user, test.findErr)
			s.e.On("ValidatePassword", "xxx", "old password").Return(test.validationErr)
			s.e.On("GeneratePassword", "new password").Return("xyz", nil)
			s.r.On("UpdateUser", s.ctx, userId, map[string]any{"password": "xyz"}).Return(test.updateErr)

			err := s.service.ChangePassword(s.ctx, userId, "new password", "old password")
			if test.expectedError != nil {
				s.Error(err)
				s.ErrorIs(err, test.expectedError)
				return
			}
			s.NoError(err)
			s.r.AssertExpectations(s.T())
		})
	}
}

func (s *serviceTestSuite) TestLogin() {
	id, err := uuid.NewV4()
	s.Require().NoError(err)
	someErr := errors.New("some")
	someUser := shared.User{Name: "name", ID: id, Password: "user password"}
	somePerm := []shared.UserAppPermission{{ID: 1, UserId: id, App: "app"}}
	tests := map[string]struct {
		user        shared.User
		findUserErr error
		validateErr error
		perm        []shared.UserAppPermission
		permErr     error
		expectError error
	}{
		"ok":               {user: someUser, perm: somePerm},
		"user not exists":  {findUserErr: someErr, expectError: someErr},
		"invalid password": {user: someUser, validateErr: someErr, expectError: someErr},
	}
	for name, test := range tests {
		s.Run(name, func() {
			s.r.On("FindUserByName", s.ctx, "name").Return(test.user, test.findUserErr)
			s.e.On("ValidatePassword", test.user.Password, "password").Return(test.validateErr)
			s.r.On("FindPermissions", s.ctx, test.user.ID).Return(test.perm, test.permErr)

			user, perm, err := s.service.Login(s.ctx, "name", "password")
			if test.expectError != nil {
				s.Error(err)
			} else {
				s.NoError(err)
				s.r.AssertExpectations(s.T())
				s.Equal(someUser, user)
				s.Equal(somePerm, perm)
			}
		})
	}
}
