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
			Date: time.Now().Add(time.Hour),
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
	var firstDiary = diaries[0] // First is latest
	assert.Equal(t, firstDiary.Content, "Ingredient")
	assert.Equal(t, firstDiary.Date.Day(), time.Now().Add(time.Hour).Day())
	assert.Equal(t, firstDiary.Severity, "cooked")
	var secondDiary = diaries[1] // Second happened earlier
	assert.Equal(t, secondDiary.Content, "Symptom")
	assert.Equal(t, secondDiary.Severity, "4")
}
