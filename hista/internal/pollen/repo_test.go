package pollen

import (
	"context"
	"testing"
	"time"

	"encore.app/errors"
	"encore.dev/et"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PollenRepoTestSuite struct {
	suite.Suite
	ctx    context.Context
	rootDb *gorm.DB
	tx     *gorm.DB
	p      *PollenRepo
}

func (s *PollenRepoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	sqlDb, err := et.NewTestDatabase(s.ctx, "hista_db")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}), &gorm.Config{TranslateError: true})
	s.rootDb = db
	if err != nil {
		panic(err)
	}
}

func (s *PollenRepoTestSuite) SetupTest() {
	s.tx = s.rootDb.Begin()
	s.Require().NoError(s.tx.Error)
	s.p = NewPollenRepo(s.tx)
}

func (s *PollenRepoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func (s *PollenRepoTestSuite) Debug() {
	s.tx = s.tx.Debug()
}

func TestPollen(t *testing.T) {
	suite.Run(t, new(PollenRepoTestSuite))
}

func (s *PollenRepoTestSuite) TestCreateAndListPollen() {
	event := &PollenEvent{CreatedAt: time.Now(), Pollens: []Pollen{
		{Type: Ambrosia, Intensity: HighPollen},
		{Type: Beifuss, Intensity: NoPollen},
	}}

	err := s.p.CreatePollen(s.ctx, event)
	s.NoError(err)

	events, err := s.p.ListPollen(s.ctx)
	s.NoError(err)
	s.Len(events, 1)
	s.Len(events[0].Pollens, 2)
}

func (s *PollenRepoTestSuite) TestDoesExistAfter() {
	createdAt := time.Date(2022, 11, 3, 13, 43, 12, 0, time.UTC)
	event := &PollenEvent{CreatedAt: createdAt}
	err := s.p.CreatePollen(s.ctx, event)
	s.Require().NoError(err)

	err = s.p.DoesExistAfter(s.ctx, createdAt.Add(time.Hour))
	s.NoError(err)

	err = s.p.DoesExistAfter(s.ctx, createdAt.Add(-time.Hour))
	s.ErrorIs(err, errors.ErrorAlreadyExists)
}

func (s *PollenRepoTestSuite) TestListPollenWithSeverity() {
	event := &PollenEvent{CreatedAt: time.Now(), Pollens: []Pollen{
		{Type: Ambrosia, Intensity: HighPollen},
		{Type: Beifuss, Intensity: NoPollen},
	}}

	err := s.p.CreatePollen(s.ctx, event)
	s.NoError(err)

	events, err := s.p.ListPollenWithSeverity(s.ctx, 0)
	s.NoError(err)
	s.Len(events, 1)
	s.Len(events[0].Pollens, 2)

	events, err = s.p.ListPollenWithSeverity(s.ctx, 1)
	s.NoError(err)
	s.Len(events, 1)
	s.Len(events[0].Pollens, 1)

	events, err = s.p.ListPollenWithSeverity(s.ctx, 7)
	s.NoError(err)
	s.Len(events, 1)
	s.Len(events[0].Pollens, 0)
}
