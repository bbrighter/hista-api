package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToMealMetaResponse(t *testing.T) {
	t.Parallel()

	var meal = Meal{
		ID:    1,
		Date:  time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC),
		Foods: []Food{{ID: 1}},
	}
	var resp MealMetaResponse = meal.ToMealMetaResponse()
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC), resp.Date)
}

func TestToMealResponse(t *testing.T) {
	t.Parallel()

	var meal = Meal{
		ID:          1,
		Date:        time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC),
		Freshness:   Fresh,
		StressLevel: 2,
		IsAlone:     true,
		Foods: []Food{{
			ID:           1,
			Ingredient:   Ingredient{ID: 10, Name: "Ingredient"},
			IngredientID: 10,
			Condition:    Cooked,
			MealID:       2,
		}},
	}

	var resp MealResponse = meal.ToMealResponse()

	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC), resp.Date)
	assert.Len(t, resp.Foods, 1)
	assert.Equal(t, Fresh, resp.Freshness)
	assert.True(t, resp.IsAlone)
	assert.Equal(t, uint8(2), resp.StressLevel)
	food := resp.Foods[0]
	assert.EqualValues(t, 1, food.ID)
	assert.Equal(t, Cooked, food.Condition)
	assert.EqualValues(t, 10, food.Ingredient.ID)
	assert.Equal(t, "Ingredient", food.Ingredient.Name)
}
