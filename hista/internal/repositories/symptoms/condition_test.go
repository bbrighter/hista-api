package symptoms

import (
	"context"
	"testing"

	"encore.app/hista/entity"
	"encore.app/shared/contextKeys"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var GUID = uuid.FromStringOrNil("cf0d4408-8db5-4572-b5d9-4ed873d1341f")

func initTest(t *testing.T) (*SymptomsRepo, context.Context) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&entity.Symptom{}, &entity.ConditionEvent{}, &entity.Condition{}, &entity.SymptomCategory{})
	assert.NoError(t, err)

	ctx := context.WithValue(t.Context(), contextKeys.Piid, GUID)
	return &SymptomsRepo{db: db}, ctx
}

func (s *SymptomRepoTestSuite) TestListConditions() {
	conditions, err := s.repo.ListConditions(s.ctx, 1)
	s.NoError(err)
	s.Len(conditions, 0)

	eventId, symptomId, _, _ := s.createTestCondition()

	conditions, err = s.repo.ListConditions(s.ctx, eventId)
	s.NoError(err)
	s.Len(conditions, 1)
	s.Equal(conditions[0].Symptom.ID, symptomId)
}

func (s *SymptomRepoTestSuite) TestCreateConditionBySymptomName() {
	eventId, _, _, catId := s.createTestCondition()
	var err error

	_, err = s.repo.CreateConditionBySymptomName(s.ctx, eventId, "name", catId)
	s.NoError(err)

	_, err = s.repo.CreateConditionBySymptomName(s.ctx, 1000, "name", catId)
	s.ErrorIs(err, gorm.ErrRecordNotFound)

	_, err = s.repo.CreateConditionBySymptomName(s.ctx, eventId, "name", 99)
	s.ErrorIs(err, gorm.ErrForeignKeyViolated)
}

func (s *SymptomRepoTestSuite) TestCreateConditionBySymptomID() {

	eventId, symptomId, _, _ := s.createTestCondition()
	_, err := s.repo.CreateConditionBySymptomID(s.ctx, eventId, symptomId)
	s.NoError(err)

	_, err = s.repo.CreateConditionBySymptomID(s.ctx, 99, symptomId)
	s.ErrorIs(err, gorm.ErrRecordNotFound)

	_, err = s.repo.CreateConditionBySymptomID(s.ctx, eventId, 99)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *SymptomRepoTestSuite) TestDeleteCondition() {
	_, _, conditionId, _ := s.createTestCondition()

	err := s.repo.DeleteCondition(s.ctx, conditionId)
	s.NoError(err)

	rows, _ := gorm.G[entity.Condition](s.tx).Count(s.ctx, "*")
	s.EqualValues(0, rows)
	rows, _ = gorm.G[entity.Symptom](s.tx).Count(s.ctx, "*")
	s.EqualValues(0, rows)
	rows, _ = gorm.G[entity.SymptomCategories](s.tx).Count(s.ctx, "*")
	s.EqualValues(1, rows)

	err = s.repo.DeleteCondition(s.ctx, 100)
	s.Error(err)
}

func (s *SymptomRepoTestSuite) TestDeleteConditionDoesNotDeleteOthersIfUsed() {
	eventId, symptomId, conditionId, _ := s.createTestCondition()
	_, err := s.repo.CreateConditionBySymptomID(s.ctx, eventId, symptomId)
	s.Require().NoError(err)

	err = s.repo.DeleteCondition(s.ctx, conditionId)
	s.NoError(err)

	rows, _ := gorm.G[entity.Condition](s.tx).Count(s.ctx, "*")
	s.EqualValues(1, rows)
	rows, _ = gorm.G[entity.Symptom](s.tx).Count(s.ctx, "*")
	s.EqualValues(1, rows)
}

func (s *SymptomRepoTestSuite) TestChangeSeverity() {
	var err error
	err = s.repo.ChangeSeverity(s.ctx, 1, entity.HighSeverity)
	s.ErrorIs(err, gorm.ErrRecordNotFound)

	_, _, conditionId, _ := s.createTestCondition()
	err = s.repo.ChangeSeverity(s.ctx, conditionId, entity.HighSeverity)
	s.NoError(err)

}

func (s *SymptomRepoTestSuite) TestGetCondition() {
	_, err := s.repo.GetCondition(s.ctx, 1)
	s.ErrorIs(err, gorm.ErrRecordNotFound)

	_, _, conId, _ := s.createTestCondition()
	con, err := s.repo.GetCondition(s.ctx, conId)
	s.NoError(err)
	s.Equal(conId, con.ID)
}
