package symptoms

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func initAPITest(t *testing.T) (*Service, context.Context, func(t *testing.T)) {
	var ctx context.Context = context.TODO()
	service, teardown := initTest(t)
	return service, ctx, teardown
}

func initTest(t *testing.T) (*Service, func(t *testing.T)) {
	service, err := initService()
	assert.NoError(t, err)

	return service, service.teardown
}

func (service Service) teardown(t *testing.T) {
	var models = []interface{}{&Condition{}, &ConditionEvent{}, &Symptom{}, &SymptomCategory{}}
	var err error
	for _, model := range models {
		err = service.db.Where("1=1").Delete(model).Error
		assert.NoError(t, err)
	}
}

func (service Service) testCreateConditionEvent(t *testing.T) ConditionEvent {
	var category = SymptomCategory{
		ID:   1,
		Name: "Category",
	}
	var err error = service.db.Create(&category).Error
	assert.NoError(t, err)
	var event = ConditionEvent{
		ID:   1,
		Date: time.Now(),
		Conditions: []Condition{
			{
				ID:               1,
				Symptom:          Symptom{ID: 1, Name: "Name", SymptomCategoryID: 1},
				SymptomID:        1,
				Severity:         1,
				ConditionEventID: 1,
			},
		},
	}
	err = service.db.Create(&event).Error
	assert.NoError(t, err)
	return event
}
