package entity

import (
	"testing"
	"time"

	"encore.app/pollen"
	"github.com/stretchr/testify/assert"
)

func TestCreateRawDiary(t *testing.T) {
	t.Parallel()

	now := time.Now()

	meals := Meals{Meal{
		Date: now,
		Foods: []Food{{
			Condition: Cooked,
			Ingredient: Ingredient{
				Name: "Ingredient",
			}}}}}
	events := ConditionEvents{ConditionEvent{
		ID:   1,
		Date: now.Add(-time.Second),
		Conditions: []Condition{{
			Severity: High,
			Symptom: Symptom{
				Name:              "Symptom",
				SymptomCategoryID: 1,
			}}}}}
	cats := SymptomCategories{SymptomCategory{
		ID:   1,
		Name: "Category",
	}}
	notes := Notes{Note{
		Text: "Note",
		Date: now.Add(-2 * time.Second)}}
	pollens := pollen.PollenEvents{pollen.PollenEvent{
		CreatedAt: now.Add(-3 * time.Second),
		Pollens: pollen.Pollens{pollen.Pollen{
			Type:      pollen.Ambrosia,
			Intensity: pollen.Medium,
		}}}}
	var diaries []RawDiary = CreateRawDiary(meals, events, cats, notes, pollens)
	assert.Len(t, diaries, 4)
	var firstDiary RawDiary = diaries[0] // First is latest
	assert.Equal(t, firstDiary.Content, "Ingredient")
	assert.Equal(t, firstDiary.Severity, "cooked")
	var secondDiary RawDiary = diaries[1] // Second happened earlier
	assert.Equal(t, secondDiary.Content, "Symptom")
	assert.Equal(t, secondDiary.Severity, "4")
	assert.Equal(t, secondDiary.Category, "Category")
	var thirdDiary RawDiary = diaries[2]
	assert.Equal(t, "Note", thirdDiary.Content)
	var fourthDiary RawDiary = diaries[3]
	assert.Equal(t, PollenType, fourthDiary.Type)
	assert.Equal(t, "Ambrosia", fourthDiary.Category)
	assert.Equal(t, "Mittlere", fourthDiary.Severity)
}
