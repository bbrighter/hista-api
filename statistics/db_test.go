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

func testInput() Input {
	var meals = meals.Meals{
		meals.Meal{
			ID:   1,
			Date: time.Now().Add(2 * time.Hour),
			Foods: []meals.Food{{
				ID: 10,
				Ingredient: meals.Ingredient{
					ID:   100,
					Name: "Ingredient",
				},
				IngredientID: 100,
				Condition:    meals.Cooked,
				MealID:       1,
			}},
		},
	}
	var events = symptoms.ConditionEvents{
		symptoms.ConditionEvent{
			ID:   1,
			Date: time.Now(),
			Conditions: []symptoms.Condition{{
				ID: 10,
				Symptom: symptoms.Symptom{
					ID:                100,
					Name:              "Symptom",
					SymptomCategoryID: 1000,
				},
				SymptomID:        100,
				Severity:         symptoms.High,
				ConditionEventID: 1,
			}},
		},
	}
	var cats = symptoms.SymptomCategories{
		symptoms.SymptomCategory{
			ID:   1000,
			Name: "Category",
		}}
	return Input{
		Meals:      meals,
		Events:     events,
		Categories: cats,
	}
}
