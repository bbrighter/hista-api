package meta

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/internal/headaches"
	"encore.app/hista/internal/status"
	"encore.app/shared/contextKeys"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repoTestSuite struct {
	suite.Suite
	ctx  context.Context
	db   *gorm.DB
	tx   *gorm.DB
	piid uuid.UUID
	m    *MetaRepo
}

const GUID_STR = "cf0d4408-8db5-4572-b5d9-4ed873d1341f"

func (s *repoTestSuite) SetupSuite() {
	s.piid = uuid.FromStringOrNil(GUID_STR)
	s.ctx = context.WithValue(context.Background(), contextKeys.Piid, s.piid)
	sqlDb, err := et.NewTestDatabase(s.ctx, "hista_db")
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

func (s *repoTestSuite) Debug() {
	s.tx = s.tx.Debug()
}

func (s *repoTestSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.Require().NoError(s.tx.Error)
	s.m = NewMetaRepo(s.tx)
}

func (s *repoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestMeta(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

func (s *repoTestSuite) TestAreTodaysEntitiesFilled_bothTrue() {
	now := time.Now()
	yesterday := time.Now().Add(-time.Hour * 24)
	tests := map[string]struct {
		createHeadacheAt       time.Time
		createStatusAt         time.Time
		expectedHeadacheFilled bool
		expectedStatusFilled   bool
	}{
		"both true":                   {createHeadacheAt: now, createStatusAt: now, expectedHeadacheFilled: true, expectedStatusFilled: true},
		"both yesterday, hence false": {createHeadacheAt: yesterday, createStatusAt: yesterday},
		"headache true, status not":   {createHeadacheAt: now, createStatusAt: yesterday, expectedHeadacheFilled: true},
	}

	for name, test := range tests {
		s.Run(name, func() {
			gorm.G[headaches.Headache](s.db).Create(s.ctx, &headaches.Headache{Date: test.createHeadacheAt, PIID: s.piid})
			gorm.G[status.Status](s.db).Create(s.ctx, &status.Status{Date: test.createStatusAt, PIID: s.piid})

			entities, err := s.m.AreTodaysEntitiesFilled(s.ctx)
			s.NoError(err)
			s.Len(entities, 2)

			s.Equal("headaches", entities[0].Entity)
			s.Equal(test.expectedHeadacheFilled, entities[0].Filled)
			s.Equal("statuses", entities[1].Entity)
			s.Equal(test.expectedStatusFilled, entities[1].Filled)
		})
	}

}

func (s *repoTestSuite) TestAreTodaysEntitiesFilled_noneTrue() {
	gorm.G[headaches.Headache](s.db).Create(s.ctx, &headaches.Headache{Date: time.Now(), PIID: s.piid})
	gorm.G[status.Status](s.db).Create(s.ctx, &status.Status{Date: time.Now(), PIID: s.piid})

	entities, err := s.m.AreTodaysEntitiesFilled(s.ctx)
	s.NoError(err)
	s.Len(entities, 2)

	s.Equal("headaches", entities[0].Entity)
	s.True(entities[0].Filled)
	s.Equal("statuses", entities[1].Entity)
	s.True(entities[1].Filled)
}
