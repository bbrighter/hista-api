package users

import (
	"testing"

	"encore.app/errors"
	"encore.dev/beta/errs"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
)

type apiTestSuite struct{ testSuite }

func (s *apiTestSuite) SetupSuite()      { s.setup() }
func (s *apiTestSuite) TearDownTest()    { s.cleanup() }
func (s *apiTestSuite) TearDownSubTest() { s.cleanup() }
func TestApi(t *testing.T)               { suite.Run(t, new(apiTestSuite)) }

func (s *apiTestSuite) TestCreateUser() {
	_, err := s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "pw"})
	s.NoError(err)
}

func (s *apiTestSuite) TestCreateUser_twice() {
	_, err := s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "pw"})
	s.NoError(err)
	_, err = s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "pw"})

	s.Error(err)
	s.assertErrCode(err, errs.AlreadyExists)
}

func (s *apiTestSuite) TestListUsers() {
	id := s.createTestUser()

	resp, err := s.service.ListUsers(s.ctx)
	s.NoError(err)
	s.Len(resp.Users, 1)
	s.Equal(UserResponse{ID: id, Name: "name"}, resp.Users[0])
}

func (s *apiTestSuite) TestDeleteUser() {
	tests := map[string]struct {
		useExistingId bool
		expectedError errs.ErrCode
	}{
		"ok":        {useExistingId: true},
		"not found": {useExistingId: false, expectedError: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			id := s.uuid()
			if test.useExistingId {
				id = s.createTestUser()
			}

			err := s.service.DeleteUser(s.ctx, id)
			s.assertErrCode(err, test.expectedError)
		})
	}
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
			_, err := s.service.CreateUser(s.ctx, UserParams{Name: "Name", Password: "Password"})
			s.Require().NoError(err)

			_, err = s.service.GetPermissions(s.ctx, LoginParams{UserName: test.userName, Password: test.password})

			if test.expectError > 0 {
				s.assertErrCode(err, test.expectError)
				return
			}
			s.NoError(err)
		})
	}
}

func (s *apiTestSuite) TestPatchPassword() {
	tests := map[string]struct {
		useWrongPassword bool
		useWrongUser     bool
		expectedError    errs.ErrCode
	}{
		"ok":             {},
		"wrong password": {useWrongPassword: true, expectedError: errs.Unauthenticated},
		"use wrong user": {useWrongUser: true, expectedError: errs.NotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			password := "password"
			id := s.createTestUserWith("name", password)

			if test.useWrongPassword {
				password = "wrong password"
			}
			if test.useWrongUser {
				id = s.uuid()
			}

			err := s.service.PatchPassword(s.ctx, id, UserPasswordChangeParams{NewPassword: "new password", OldPassword: password})
			s.assertErrCode(err, test.expectedError)
		})
	}
}

func (s *apiTestSuite) TestGetSettings() {

	tests := map[string]struct {
		useExistingId bool
		err           error
	}{
		"ok":        {useExistingId: true},
		"not found": {useExistingId: false, err: errors.ErrorNotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			idResp, err := s.service.CreateUser(s.ctx, UserParams{})
			s.Require().NoError(err)
			userId := idResp.UserId
			if !test.useExistingId {
				userId = s.uuid()
			}

			settingsResp, err := s.service.GetUserSettings(s.ctx, userId)
			if test.err != nil {
				s.Error(err)
				s.ErrorIs(err, test.err)
				return
			}

			s.Equal("de-DE", settingsResp.Language)
			s.Equal("spinner", settingsResp.LoadingMode)
		})
	}
}

func (s *apiTestSuite) TestChangeSettings() {
	defaultLoadingMode := "spinner"
	defaultLanguage := "de-DE"
	newLoadingMode := "fancy"
	newLanguage := "en-US"

	tests := map[string]struct {
		useExistingId  bool
		setLoadingMode bool
		setLanguage    bool
		expectedErr    error
	}{
		"ok":                   {useExistingId: true, setLoadingMode: true, setLanguage: true},
		"change only language": {useExistingId: true, setLanguage: true},
		"not found":            {useExistingId: false, expectedErr: errors.ErrorNotFound},
	}
	for name, test := range tests {
		s.Run(name, func() {
			idResp, err := s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "password"})
			s.Require().NoError(err)

			otherId, err := uuid.NewV4()
			s.Require().NoError(err)

			var id = idResp.UserId
			if !test.useExistingId {
				id = otherId
			}

			params := UserSettingsPatchParams{}
			if test.setLanguage {
				params.Language = option.FromComparable(newLanguage)
			}
			if test.setLoadingMode {
				params.LoadingMode = option.FromComparable(newLoadingMode)
			}

			err = s.service.PatchUserSettings(s.ctx, id, params)

			if test.expectedErr != nil {
				s.Error(err)
				s.assertErrCode(err, errs.NotFound)
				return
			}

			s.NoError(err)
			settingsResp, err := s.service.GetUserSettings(s.ctx, idResp.UserId)
			s.NoError(err)
			if test.setLanguage {
				s.Equal(newLanguage, settingsResp.Language)
			} else {
				s.Equal(defaultLanguage, settingsResp.Language)
			}

			if test.setLoadingMode {
				s.Equal(newLoadingMode, settingsResp.LoadingMode)
			} else {
				s.Equal(defaultLoadingMode, settingsResp.LoadingMode)
			}

		})
	}

}

func (s *apiTestSuite) TestChangeSettings_OnlyAffectedChanged() {
	idResp, err := s.service.CreateUser(s.ctx, UserParams{Name: "name", Password: "password"})
	s.NoError(err)

	err = s.service.PatchUserSettings(s.ctx, idResp.UserId, UserSettingsPatchParams{
		LoadingMode: option.None[string](),
		Language:    option.Some("en-US"),
	})
	s.NoError(err)

	settingsResp, err := s.service.GetUserSettings(s.ctx, idResp.UserId)
	s.NoError(err)
	s.Equal("en-US", settingsResp.Language)
	s.Equal("spinner", settingsResp.LoadingMode)
}
