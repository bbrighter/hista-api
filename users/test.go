package users

import (
	"context"

	"encore.dev/beta/errs"
	"encore.dev/et"
	uuid "encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type testSuite struct {
	suite.Suite

	ctx     context.Context
	service *Service
	db      *gorm.DB
}

func (s *testSuite) setup() {
	s.ctx = context.Background()

	sqlDB, err := et.NewTestDatabase(s.ctx, "users_db")
	s.Require().NoError(err)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB.Stdlib(),
	}))
	s.Require().NoError(err)

	s.db = db
	s.service = setupService(db, 4)
}

func (s *testSuite) cleanup() {
	tables := []string{
		"users",
		"user_product_instances",
		"user_app_permissions",
	}

	for _, table := range tables {
		err := s.db.Exec("DELETE FROM " + table).Error
		s.Require().NoError(err)
	}
}

func (s *testSuite) assertErrCode(err error, code errs.ErrCode) {
	if code == 0 {
		return
	}
	e, ok := err.(*errs.Error)
	s.True(ok)
	s.Equal(code, e.Code)
}

func (s *testSuite) createTestUser() uuid.UUID {
	return s.createTestUserWith("name", "pw")
}

func (s *testSuite) createTestUserWith(name string, pw string) uuid.UUID {
	idResp, err := s.service.CreateUser(
		s.ctx,
		UserParams{Name: name, Password: pw},
	)
	s.Require().NoError(err)

	return idResp.UserId
}

func (s *testSuite) uuid() uuid.UUID {
	id, err := uuid.NewV4()
	s.Require().NoError(err)
	return id
}
