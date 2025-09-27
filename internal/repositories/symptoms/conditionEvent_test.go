package symptoms

import (
	"testing"
	"time"

	"encore.app/entity"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func (repo *SymptomsRepo) createTestConditionEvent(t *testing.T) uint {
	var event = entity.ConditionEvent{PIID: GUID}
	err := repo.db.Create(&event).Error
	assert.NoError(t, err)
	return event.ID
}

func TestCreateConditionEvent(t *testing.T) {
	repo, ctx := initTest(t)

	var event = &entity.ConditionEvent{Date: time.Now(), PIID: GUID}
	err := repo.CreateConditionEvent(ctx, event)
	assert.NoError(t, err)
}

func TestListConditionEvents(t *testing.T) {
	repo, ctx := initTest(t)

	events, err := repo.ListConditionEvents(ctx)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	repo.createTestConditionEvent(t)
	events, err = repo.ListConditionEvents(ctx)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
}

func TestGetConditionEvent(t *testing.T) {
	repo, ctx := initTest(t)

	_, err := repo.GetConditionEvent(ctx, 1)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	id := repo.createTestConditionEvent(t)
	event, err := repo.GetConditionEvent(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, event.ID)
}

func TestDeleteConditionEvent(t *testing.T) {
	repo, ctx := initTest(t)

	err := repo.DeleteConditionEvent(ctx, 1)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	id := repo.createTestConditionEvent(t)
	err = repo.DeleteConditionEvent(ctx, id)
	assert.NoError(t, err)
}

func TestDeleteConditionEventOfOTherPiid(t *testing.T) {
	repo, ctx := initTest(t)
	event := entity.ConditionEvent{PIID: uuid.FromStringOrNil("5f0347ae-38e8-49f4-9707-c3f4500e1768")}
	err := repo.db.Create(&event).Error
	require.NoError(t, err)

	err = repo.DeleteConditionEvent(ctx, event.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestDeleteConditionEventWithChildren(t *testing.T) {
	repo, ctx := initTest(t)

	id := repo.createTestConditionEvent(t)
	repo.db.Create(&entity.SymptomCategory{
		ID:   1,
		Name: "cat",
		PIID: GUID,
		Symptoms: []entity.Symptom{
			{
				ID:                10,
				Name:              "symptom",
				SymptomCategoryID: 1,
				PIID:              GUID,
			},
		},
	})
	repo.db.Create(&entity.Condition{ID: 100, SymptomID: 10, ConditionEventID: id, PIID: GUID})
	err := repo.DeleteConditionEvent(ctx, id)
	assert.NoError(t, err)

	var symptoms entity.Symptoms
	var conditions entity.Conditions
	assert.EqualValues(t, 0, repo.db.Find(&symptoms).RowsAffected)
	assert.EqualValues(t, 0, repo.db.Find(&conditions).RowsAffected)
}

func TestPatchConditionEvent(t *testing.T) {
	repo, ctx := initTest(t)

	err := repo.PatchConditionEvent(ctx, 1, time.Now())
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	id := repo.createTestConditionEvent(t)
	err = repo.PatchConditionEvent(ctx, id, time.Now())
	assert.NoError(t, err)
}
