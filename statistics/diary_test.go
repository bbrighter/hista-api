package statistics

import (
	"testing"
	"time"

	"encore.app/meals"
	"encore.app/symptoms"
	"github.com/stretchr/testify/assert"
)

func TestDiaryFrom(t *testing.T) {
	t.Parallel()
	var meals = meals.Meals{
		meals.Meal{
			ID:   1,
			Date: time.Now(),
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
	var diaries []RawDiary
	diaries = diaryFrom(meals, events)
	assert.Len(t, diaries, 2)
	var firstDiary = diaries[0]
	assert.Equal(t, firstDiary.Content, "Ingredient")
	assert.Equal(t, firstDiary.Date.Day(), time.Now().Day())
}

func TestGetDiaryData(t *testing.T) {
	service, teardown := initTest(t)
	service.testCreateData(t)
	defer teardown(t)

	meals, events := getDiaryData(service)
	assert.Len(t, meals, 1)
	meal := meals[0]
	assert.Len(t, meal.Foods, 1)
	assert.Equal(t, meal.Foods[0].Ingredient.Name, "Ingredient")

	assert.Len(t, events, 1)
	event := events[0]
	assert.Len(t, event.Conditions, 1)
	assert.Equal(t, event.Conditions[0].Symptom.Name, "Symptom")
}
