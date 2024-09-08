package symptoms

import (
	"testing"
	"time"

	"encore.app/entity"
	"github.com/stretchr/testify/assert"
)

func (repo *SymptomsRepo) createTestConditionEvent(t *testing.T) uint {
	var event entity.ConditionEvent
	err := repo.db.Create(&event).Error
	assert.NoError(t, err)
	return event.ID
}

func TestCreateConditionEvent(t *testing.T) {
	repo := initTest(t)

	var event = &entity.ConditionEvent{Date: time.Now()}
	err := repo.CreateConditionEvent(event)
	assert.NoError(t, err)
}

func TestListConditionEvents(t *testing.T) {
	repo := initTest(t)

	var events entity.ConditionEvents
	events = repo.ListConditionEvents()
	assert.Len(t, events, 0)

	repo.createTestConditionEvent(t)
	events = repo.ListConditionEvents()
	assert.Len(t, events, 1)
}

func TestGetConditionEvent(t *testing.T) {
	repo := initTest(t)

	_, err := repo.GetConditionEvent(1)
	assert.EqualError(t, err, "not_found: not found")

	id := repo.createTestConditionEvent(t)
	event, err := repo.GetConditionEvent(id)
	assert.NoError(t, err)
	assert.Equal(t, id, event.ID)
}

func TestDeleteConditionEvent(t *testing.T) {
	repo := initTest(t)

	err := repo.DeleteConditionEvent(1)
	assert.EqualError(t, err, "not_found: not found")

	id := repo.createTestConditionEvent(t)
	err = repo.DeleteConditionEvent(id)
	assert.NoError(t, err)
}

func TestDeleteConditionEventWithChildren(t *testing.T) {
	repo := initTest(t)

	id := repo.createTestConditionEvent(t)
	repo.db.Create(&entity.SymptomCategory{
		ID:   1,
		Name: "cat",
		Symptoms: []entity.Symptom{
			{
				ID:                10,
				Name:              "symptom",
				SymptomCategoryID: 1,
			},
		},
	})
	repo.db.Create(&entity.Condition{ID: 100, SymptomID: 10, ConditionEventID: id})
	err := repo.DeleteConditionEvent(id)
	assert.NoError(t, err)

	var symptoms entity.Symptoms
	var cats entity.SymptomCategories
	var conditions entity.Conditions
	assert.EqualValues(t, 0, repo.db.Find(&symptoms).RowsAffected)
	assert.EqualValues(t, 0, repo.db.Find(&cats).RowsAffected)
	assert.EqualValues(t, 0, repo.db.Find(&conditions).RowsAffected)
}

func TestPatchConditionEvent(t *testing.T) {
	repo := initTest(t)

	err := repo.PatchConditionEvent(1, time.Now())
	assert.EqualError(t, err, "not_found: not found")

	id := repo.createTestConditionEvent(t)
	err = repo.PatchConditionEvent(id, time.Now())
	assert.NoError(t, err)
}
