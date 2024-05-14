package statistics

import (
	"context"
	"testing"
	"time"

	"encore.app/meals"
	"encore.app/symptoms"
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
	var models = []interface{}{
		&meals.Food{},
		&meals.Meal{},
		&meals.Ingredient{},
		&symptoms.Condition{},
		&symptoms.ConditionEvent{},
		&symptoms.Symptom{},
		&symptoms.SymptomCategory{},
	}
	var err error
	for _, model := range models {
		err = service.db.Debug().Where("1=1").Delete(model).Error
		assert.NoError(t, err)
	}
}

func (service *Service) testCreateData(t *testing.T) (meals.Meal, symptoms.ConditionEvent) {
	var food = meals.Food{
		ID: 1,
		Ingredient: meals.Ingredient{
			ID:   1,
			Name: "Ingredient",
		},
		IngredientID: 1,
	}
	var meal = meals.Meal{
		ID:    1,
		Date:  time.Now(),
		Foods: []meals.Food{food},
	}
	var err error = service.db.Create(&meal).Error
	assert.NoError(t, err)

	var symptomCategory = symptoms.SymptomCategory{
		ID:   1,
		Name: "Category",
	}
	err = service.db.Create(&symptomCategory).Error
	assert.NoError(t, err)
	var event = symptoms.ConditionEvent{
		ID:   1,
		Date: time.Now(),
		Conditions: []symptoms.Condition{{
			ID:               10,
			SymptomID:        100,
			Severity:         symptoms.High,
			ConditionEventID: 1,
			Symptom: symptoms.Symptom{
				ID:                100,
				Name:              "Symptom",
				SymptomCategoryID: 1,
			}}},
	}
	err = service.db.Create(&event).Error
	assert.NoError(t, err)

	return meal, event
}
