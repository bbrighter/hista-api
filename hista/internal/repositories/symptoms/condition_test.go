package symptoms

import (
	"context"
	"testing"

	"encore.app/hista/entity"
	"encore.dev/types/uuid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var GUID = uuid.FromStringOrNil("cf0d4408-8db5-4572-b5d9-4ed873d1341f")

func initTest(t *testing.T) (*SymptomsRepo, context.Context) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&entity.Symptom{}, &entity.ConditionEvent{}, &entity.Condition{}, &entity.SymptomCategory{})
	assert.NoError(t, err)

	ctx := context.WithValue(t.Context(), "piid", GUID)
	return &SymptomsRepo{db: db}, ctx
}

func (repo *SymptomsRepo) createTestCondition(t *testing.T) (uint, uint, uint, uint) {
	var eventId uint = 1
	var conditionId uint = 10
	var symptomId uint = 100
	var catId uint = 1000
	err := repo.db.Create(&entity.ConditionEvent{ID: eventId, PIID: GUID}).Error
	require.NoError(t, err)
	err = repo.db.Create(&entity.SymptomCategory{ID: catId, PIID: GUID}).Error
	require.NoError(t, err)
	err = repo.db.Create(&entity.Symptom{ID: symptomId, SymptomCategoryID: catId, PIID: GUID, SymptomCategoryPIID: GUID}).Error
	require.NoError(t, err)
	err = repo.db.Create(&entity.Condition{
		ID:                 conditionId,
		SymptomID:          symptomId,
		ConditionEventID:   eventId,
		PIID:               GUID,
		SymptomPIID:        GUID,
		ConditionEventPIID: GUID,
	}).Error
	require.NoError(t, err)

	return eventId, symptomId, conditionId, catId
}

func TestListConditions(t *testing.T) {
	repo, ctx := initTest(t)

	conditions, err := repo.ListConditions(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, conditions, 0)

	eventId, symptomId, _, _ := repo.createTestCondition(t)

	conditions, err = repo.ListConditions(ctx, eventId)
	assert.NoError(t, err)
	assert.Len(t, conditions, 1)
	assert.Equal(t, conditions[0].Symptom.ID, symptomId)
}

func TestCreateConditionBySymptomName(t *testing.T) {
	repo, ctx := initTest(t)

	eventId, _, oldConditionId, catId := repo.createTestCondition(t)
	id, err := repo.CreateConditionBySymptomName(ctx, eventId, "name", catId)
	assert.NoError(t, err)
	assert.EqualValues(t, oldConditionId+1, id)

	_, err = repo.CreateConditionBySymptomName(ctx, 1000, "name", catId)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	id, err = repo.CreateConditionBySymptomName(ctx, eventId, "name", 99)
	assert.NoError(t, err)
	assert.EqualValues(t, oldConditionId+2, id)
}

func TestCreateConditionBySymptomID(t *testing.T) {
	repo, ctx := initTest(t)

	eventId, symptomId, conditionId, _ := repo.createTestCondition(t)
	id, err := repo.CreateConditionBySymptomID(ctx, eventId, symptomId)
	assert.NoError(t, err)
	assert.Equal(t, conditionId+1, id)

	_, err = repo.CreateConditionBySymptomID(ctx, 99, symptomId)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	_, err = repo.CreateConditionBySymptomID(ctx, eventId, 99)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

}

func TestDeleteCondition(t *testing.T) {
	repo, ctx := initTest(t)

	_, _, conditionId, _ := repo.createTestCondition(t)

	err := repo.DeleteCondition(ctx, conditionId)
	assert.NoError(t, err)
	var conditions []entity.Condition
	repo.db.Find(&conditions)
	assert.Len(t, conditions, 0)
	var symptoms entity.Symptoms
	repo.db.Find(&symptoms)
	assert.Len(t, symptoms, 0)
	var cats entity.SymptomCategories
	repo.db.Find(&cats)
	assert.Len(t, cats, 1, "categories are not deleted, even if empty")
}

func TestChangeSeverity(t *testing.T) {
	repo, ctx := initTest(t)

	var err error
	err = repo.ChangeSeverity(ctx, 1, entity.HighSeverity)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	_, _, conditionId, _ := repo.createTestCondition(t)
	err = repo.ChangeSeverity(ctx, conditionId, entity.HighSeverity)
	assert.NoError(t, err)

}

func TestGetCondition(t *testing.T) {
	repo, ctx := initTest(t)

	_, err := repo.GetCondition(ctx, 1)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	_, _, conId, _ := repo.createTestCondition(t)
	con, err := repo.GetCondition(ctx, conId)
	assert.NoError(t, err)
	assert.Equal(t, conId, con.ID)
}
