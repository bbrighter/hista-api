package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateRawDiary(t *testing.T) {
	t.Parallel()

	now := time.Now()

	meals := Meals{&Meal{
		Date: now,
		Foods: []Food{{
			Condition: Cooked,
			Ingredient: Ingredient{
				Name: "Ingredient",
			}}}}}
	events := ConditionEvents{&ConditionEvent{
		ID:   1,
		Date: now.Add(-time.Second),
		Conditions: []Condition{{
			Severity: HighSeverity,
			Symptom: Symptom{
				Name:              "Symptom",
				SymptomCategoryID: 1,
			}}}}}
	cats := SymptomCategories{&SymptomCategory{
		ID:   1,
		Name: "Category",
	}}
	notes := Notes{&Note{
		Text: "Note",
		Date: now.Add(-2 * time.Second)}}
	pollens := PollenEvents{PollenEvent{
		CreatedAt: now.Add(-3 * time.Second),
		Pollens: Pollens{Pollen{
			Type:      Ambrosia,
			Intensity: MediumPollen,
		}}}}
	intakes := []*Intake{
		{ID: 1, Date: now.Add(-10 * time.Second), Medicine: Medicine{ID: 1, Name: "Medicine"}},
	}
	var diaries []RawDiary = CreateRawDiary(meals, events, cats, notes, pollens, intakes)
	assert.Len(t, diaries, 5)
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
	assert.Equal(t, DiaryPollen, fourthDiary.Type)
	assert.Equal(t, "Ambrosia", fourthDiary.Category)
	assert.Equal(t, "Mittlere", fourthDiary.Severity)
	var fithDiary RawDiary = diaries[4]
	assert.Equal(t, DiaryIntake, fithDiary.Type)
	assert.Equal(t, "Medicine", fithDiary.Content)
}
