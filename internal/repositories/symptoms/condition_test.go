package symptoms

import (
	"testing"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func initTest(t *testing.T) *SymptomsRepo {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&entity.Symptom{}, &entity.ConditionEvent{}, &entity.Condition{}, &entity.SymptomCategory{})
	assert.NoError(t, err)
	return &SymptomsRepo{db: db}
}

func (repo *SymptomsRepo) createTestCondition(t *testing.T) (uint, uint, uint, uint) {
	var eventId uint = 1
	var conditionId uint = 10
	var symptomId uint = 100
	var catId uint = 1000
	err := repo.db.Create(&entity.ConditionEvent{ID: eventId}).Error
	assert.NoError(t, err)
	err = repo.db.Create(&entity.SymptomCategory{ID: catId}).Error
	assert.NoError(t, err)
	err = repo.db.Create(&entity.Symptom{ID: symptomId, SymptomCategoryID: catId}).Error
	assert.NoError(t, err)
	err = repo.db.Create(&entity.Condition{
		ID:               conditionId,
		SymptomID:        symptomId,
		ConditionEventID: eventId}).Error
	assert.NoError(t, err)

	return eventId, symptomId, conditionId, catId
}

func TestListConditions(t *testing.T) {
	repo := initTest(t)

	var conditions entity.Conditions
	conditions = repo.ListConditions(1)
	assert.Len(t, conditions, 0)

	eventId, symptomId, _, _ := repo.createTestCondition(t)

	conditions = repo.ListConditions(eventId)
	assert.Len(t, conditions, 1)
	assert.Equal(t, conditions[0].Symptom.ID, symptomId)
}

func TestCreateConditionBySymptomName(t *testing.T) {
	repo := initTest(t)

	eventId, _, oldConditionId, catId := repo.createTestCondition(t)
	conditionId, err := repo.CreateConditionBySymptomName(eventId, "name", catId)
	assert.NoError(t, err)
	assert.Equal(t, oldConditionId+1, conditionId)

	_, err = repo.CreateConditionBySymptomName(99, "name", catId)
	assert.EqualError(t, err, "not_found: not found")

	conditionId, err = repo.CreateConditionBySymptomName(eventId, "name", 99)
	assert.NoError(t, err)
	assert.Equal(t, oldConditionId+2, conditionId)
}

func TestCreateConditionBySymptomID(t *testing.T) {
	repo := initTest(t)

	eventId, symptomId, conditionId, _ := repo.createTestCondition(t)
	id, err := repo.CreateConditionBySymptomID(eventId, symptomId)
	assert.NoError(t, err)
	assert.Equal(t, conditionId+1, id)

	_, err = repo.CreateConditionBySymptomID(99, symptomId)
	assert.EqualError(t, err, "not_found: not found")

	_, err = repo.CreateConditionBySymptomID(eventId, 99)
	assert.EqualError(t, err, "not_found: not found")

}

func TestDeleteCondition(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.db.Preload(clause.Associations).Find(&entity.Condition{}, &entity.Condition{SymptomID: 1}).Error
	assert.NoError(t, err)

	err = repo.DeleteCondition(1)
	assert.EqualError(t, err, "not_found: not found")

	_, _, conditionId, _ := repo.createTestCondition(t)
	err = repo.DeleteCondition(conditionId)
	assert.NoError(t, err)

	var conditions entity.Conditions
	repo.db.Find(&conditions)
	assert.Len(t, conditions, 0)
	var symptoms entity.Symptoms
	repo.db.Find(&symptoms)
	assert.Len(t, symptoms, 0)
	var cats entity.SymptomCategories
	repo.db.Find(&cats)
	assert.Len(t, cats, 0, "len of cats")
}

func TestChangeSeverity(t *testing.T) {
	repo := initTest(t)

	var err error
	err = repo.ChangeSeverity(1, entity.HighSeverity)
	assert.EqualError(t, err, "not_found: not found")

	_, _, conditionId, _ := repo.createTestCondition(t)
	err = repo.ChangeSeverity(conditionId, entity.HighSeverity)
	assert.NoError(t, err)

}
