package headaches

import (
	"context"
	"testing"
	"time"

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
	h    *HeadacheRepo
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
	s.h = NewHeadacheRepo(s.tx)
}

func (s *repoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestHeadaches(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

func (s *repoTestSuite) createHeadache() Headache {
	var headache = &Headache{
		Date:        time.Now(),
		Severity:    3,
		Types:       HeadacheTypes{Dull},
		Positions:   HeadachePositions{Back, Ear},
		Symptoms:    HeadacheSymptoms{},
		Description: "description",
	}
	err := s.h.CreateHeadache(s.ctx, headache)
	s.Require().NoError(err)
	return *headache
}

func (s *repoTestSuite) TestListHeadaches() {
	headaches, err := s.h.ListHeadaches(s.ctx)
	s.NoError(err)
	s.Len(headaches, 0)

	s.createHeadache()
	headaches, err = s.h.ListHeadaches(s.ctx)
	s.NoError(err)
	s.Len(headaches, 1)
}

func (s *repoTestSuite) TestFirstHeadache() {
	id := s.createHeadache().ID

	headache, err := s.h.FirstHeadache(s.ctx, id)
	s.NoError(err)
	s.EqualValues(3, headache.Severity)
}

func (s *repoTestSuite) TestFirstHeadacheNotFound() {
	_, err := s.h.FirstHeadache(s.ctx, 100)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *repoTestSuite) TestUpdateHeadache() {
	id := s.createHeadache().ID

	values := map[string]any{"severity": 1}
	err := s.h.UpdateHeadache(s.ctx, id, values)
	s.NoError(err)
	headache, _ := s.h.FirstHeadache(s.ctx, id)
	s.EqualValues(1, headache.Severity)

	values2 := map[string]any{"types": HeadacheTypes{Dull, Pulsating}}
	err = s.h.UpdateHeadache(s.ctx, id, values2)
	s.NoError(err)
	headache, _ = s.h.FirstHeadache(s.ctx, id)
	s.EqualValues(HeadacheTypes{Dull, Pulsating}, headache.Types)
}
