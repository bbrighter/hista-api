package statistics

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"encore.app/entity"
	"encore.app/notes"
	"encore.app/pollen"
	"encore.app/symptoms"
	"github.com/stretchr/testify/assert"
)

var testMeal *entity.Meal
var testIngredient *entity.Ingredient
var testEvent *symptoms.ConditionEvent
var testSymptom *symptoms.Symptom

func initAPITest(t *testing.T) (*Service, context.Context) {
	var ctx context.Context = context.TODO()
	service := initTest(t)
	return service, ctx
}

func initTest(t *testing.T) *Service {
	service, err := initService()
	assert.NoError(t, err)
	service.initData()

	return service
}

func (service *Service) initData() {
	var meal = entity.Meal{Date: time.Date(2021, 1, 1, 1, 0, 0, 0, time.Local)}
	service.db.FirstOrCreate(&meal, &meal)
	var ingredient = entity.Ingredient{Name: "statistics_ingredient"}
	service.db.FirstOrCreate(&ingredient, &ingredient)
	var food = entity.Food{IngredientID: ingredient.ID, MealID: meal.ID, Condition: entity.Cooked}
	service.db.FirstOrCreate(&food, &food)

	var event = symptoms.ConditionEvent{Date: time.Date(2021, 1, 1, 3, 0, 0, 0, time.Local)}
	service.db.FirstOrCreate(&event, &event)
	var symptomCategory = symptoms.SymptomCategory{Name: "statistics_category"}
	service.db.FirstOrCreate(&symptomCategory, &symptomCategory)
	var symptom = symptoms.Symptom{Name: "statistics_symptom", SymptomCategoryID: symptomCategory.ID}
	service.db.FirstOrCreate(&symptom, &symptom)
	var condition = symptoms.Condition{SymptomID: symptom.ID, ConditionEventID: event.ID, Severity: symptoms.High}
	service.db.FirstOrCreate(&condition, &condition)

	testMeal = &meal
	testEvent = &event
	testSymptom = &symptom
	testIngredient = &ingredient
}

func testInput() Input {
	var meals = entity.Meals{
		entity.Meal{
			ID:   1,
			Date: time.Now().Add(2 * time.Hour),
			Foods: []entity.Food{{
				ID: 10,
				Ingredient: entity.Ingredient{
					ID:   100,
					Name: "Ingredient",
				},
				IngredientID: 100,
				Condition:    entity.Cooked,
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
	var notes = notes.Notes{
		notes.Note{
			ID:   5,
			Date: time.Now().Add(-time.Hour),
			Text: "Note",
		}}
	var pollens = []pollen.PollenEvent{{
		ID:        7,
		CreatedAt: time.Now().Add(-time.Hour * 2),
		Pollens: pollen.Pollens{
			pollen.Pollen{
				ID:            70,
				PollenEventID: 7,
				Type:          pollen.Ambrosia,
				Intensity:     pollen.Medium,
			},
		},
	}}
	return Input{
		Meals:      meals,
		Events:     events,
		Categories: cats,
		Notes:      notes,
		Pollens:    pollens,
	}
}
