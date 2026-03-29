package symptoms

import (
	"context"
	"testing"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SymptomRepoTestSuite struct {
	suite.Suite
	ctx  context.Context
	db   *gorm.DB
	tx   *gorm.DB
	repo *SymptomsRepo
	piid uuid.UUID
}

func (s *SymptomRepoTestSuite) SetupSuite() {
	s.piid = uuid.FromStringOrNil("cf0d4408-8db5-4572-b5d9-4ed873d1341f")
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

func (s *SymptomRepoTestSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.Require().NoError(s.tx.Error)
	s.repo = NewSymptomsRepo(s.tx)
}

func (s *SymptomRepoTestSuite) TearDownTest() {
	err := s.tx.Rollback().Error
	s.Require().NoError(err)
}

func TestSymptomRepoTestSuite(t *testing.T) {
	suite.Run(t, new(SymptomRepoTestSuite))
}

func (s *SymptomRepoTestSuite) createTestCondition() (uint, uint, uint, uint) {
	var eventId uint = 1
	var conditionId uint = 10
	var symptomId uint = 100
	var catId uint = 1000
	err := s.repo.db.Create(&entity.ConditionEvent{ID: eventId, PIID: GUID}).Error
	s.Require().NoError(err)
	err = s.repo.db.Create(&entity.SymptomCategory{ID: catId, PIID: GUID}).Error
	s.Require().NoError(err)
	err = s.repo.db.Create(&entity.Symptom{ID: symptomId, SymptomCategoryID: catId, PIID: GUID, SymptomCategoryPIID: GUID}).Error
	s.Require().NoError(err)
	err = s.repo.db.Create(&entity.Condition{
		ID:                 conditionId,
		SymptomID:          symptomId,
		ConditionEventID:   eventId,
		PIID:               GUID,
		SymptomPIID:        GUID,
		ConditionEventPIID: GUID,
	}).Error
	s.Require().NoError(err)

	return eventId, symptomId, conditionId, catId
}
