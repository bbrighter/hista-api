package statistics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiaryFrom(t *testing.T) {
	t.Parallel()

	var input Input = testInput()
	var diaries []RawDiary = diaryFrom(input)
	assert.Len(t, diaries, 2)
	var firstDiary RawDiary = diaries[0] // First is latest
	assert.Equal(t, firstDiary.Content, "Ingredient")
	assert.Equal(t, firstDiary.Date.Day(), time.Now().Add(2*time.Hour).Day())
	assert.Equal(t, firstDiary.Severity, "cooked")
	var secondDiary RawDiary = diaries[1] // Second happened earlier
	assert.Equal(t, secondDiary.Content, "Symptom")
	assert.Equal(t, secondDiary.Severity, "4")
	assert.Equal(t, secondDiary.Category, "Category")
}

func TestMealBasedDiary(t *testing.T) {
	t.Parallel()

	var input Input = testInput()
	var diaries []MealBasedDiary = mealBasedDiary(input)
	assert.Len(t, diaries, 1)
	var diary MealBasedDiary = diaries[0]
	assert.Equal(t, diary.Food, "Ingredient")
	assert.Equal(t, diary.Condition, "cooked")
	assert.Len(t, diary.SymptomsWithin1h, 0)
	assert.Len(t, diary.SymptomsWithin12h, 1)
	assert.Len(t, diary.SymptomsWithin24h, 1)
}
