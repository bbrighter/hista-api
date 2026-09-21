package users

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type integrationTestSuite struct{ testSuite }

func (s *integrationTestSuite) SetupSuite()   { s.testSuite.setup() }
func (s *integrationTestSuite) TearDownTest() { s.testSuite.cleanup() }
func TestIntegration(t *testing.T)            { suite.Run(t, new(integrationTestSuite)) }

func (s *integrationTestSuite) TestUser() {
	userIdResp, err := s.service.CreateUser(s.ctx, UserParams{Name: "Name", Password: "Password"})
	s.NoError(err)
	userId := userIdResp.UserId

	usersResp, err := s.service.ListUsers(s.ctx)
	s.NoError(err)
	s.Len(usersResp.Users, 1)
	s.Equal("Name", usersResp.Users[0].Name)
	s.Equal(userId, usersResp.Users[0].ID)

	sameUserId, err := s.service.Exists(s.ctx, "Name")
	s.NoError(err)
	s.Equal(userId, sameUserId.UserId)

	err = s.service.PatchPassword(s.ctx, userId, UserPasswordChangeParams{NewPassword: "New password", OldPassword: "Password"})
	s.NoError(err)

	settingsResp, err := s.service.GetUserSettings(s.ctx, userId)
	s.NoError(err)
	s.Equal("spinner", settingsResp.LoadingMode)
	s.Equal("de-DE", settingsResp.Language)

	err = s.service.DeleteUser(s.ctx, userId)
	s.NoError(err)

	usersResp, err = s.service.ListUsers(s.ctx)
	s.NoError(err)
	s.Len(usersResp.Users, 0)
}
