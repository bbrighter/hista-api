package statistics

import (
	"context"
	_ "embed"
	"log"
	"testing"
	"time"

	"encore.app/meals"
	"encore.app/symptoms"
	"encore.dev"
	"github.com/stretchr/testify/assert"
)

func initAPITest(t *testing.T) (*Service, context.Context) {
	var ctx context.Context = context.TODO()
	service := initTest(t)
	return service, ctx
}

//go:embed fixtures.sql
var fixtures string

func initTest(t *testing.T) *Service {
	service, err := initService()
	if encore.Meta().Environment.Cloud == encore.CloudLocal {
		if _, err := histaDB.Exec(context.Background(), fixtures); err != nil {
			log.Fatalln("unable to add fixtures:", err)
		}
	}
	assert.NoError(t, err)

	return service
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
