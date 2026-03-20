package pollen

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.dev/et"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PollenRepoTestSuite struct {
	suite.Suite
	repo *PollenRepo
	ctx  context.Context
	db   *gorm.DB
	tx   *gorm.DB
}

func (s *PollenRepoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	sqlDb, err := et.NewTestDatabase(s.ctx, "hista_db")
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}
	s.db = db
}

func (s *PollenRepoTestSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.Require().NoError(s.tx.Error)
	s.repo = NewPollenRepo(s.tx)
}

func (s *PollenRepoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestPollenRepo(t *testing.T) {
	suite.Run(t, new(PollenRepoTestSuite))
}

func (s *PollenRepoTestSuite) createEvent(time time.Time) {
	event := &entity.PollenEvent{
		CreatedAt: time,
		Pollens: entity.Pollens{
			{Type: entity.Ambrosia, Intensity: entity.HighPollen},
			{Type: entity.Beifuss, Intensity: entity.NoPollen},
		},
	}
	err := gorm.G[entity.PollenEvent](s.tx).Create(s.ctx, event)
	s.Require().NoError(err)
}
