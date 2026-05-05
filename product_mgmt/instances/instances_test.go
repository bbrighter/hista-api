package instances

import (
	"context"
	"testing"

	"encore.dev/et"
	"github.com/stretchr/testify/suite"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const productTestId string = "prod"

type TestSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	inst *InstanceRepo
	ctx  context.Context
}

func (s *TestSuite) SetupSuite() {
	s.ctx = context.Background()
	sqlDb, err := et.NewTestDatabase(s.ctx, "product_mgmt_db")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	s.db = db
	if err != nil {
		panic(err)
	}
}

func (s *TestSuite) SetupSubTest() {
	s.tx = s.db.Begin()
	s.Require().NoError(s.tx.Error)
	s.inst = NewInstanceRepo(s.tx)
}

func (s *TestSuite) TeardownSubTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestCreateInstance() {
	tests := map[string]struct {
		productId     string
		expectedError bool
	}{
		"ok": {productId: productTestId},
	}
	for name, test := range tests {
		s.Run(name, func() {
			instance, err := s.inst.Create(s.ctx, "instance name", test.productId)
			if test.expectedError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal("instance name", instance.Name)
				s.NotNil(instance.ID)
				instances, err := gorm.G[Instance](s.tx).Find(s.ctx)
				s.NoError(err)
				s.Len(instances, 1)
			}
		})
	}
}
