package diary

import (
	"testing"
	"time"

	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/medicines"
	"encore.app/hista/internal/notes"
	"encore.app/hista/internal/pollen"
	"encore.app/hista/internal/symptoms"
	"github.com/stretchr/testify/assert"
)

func TestMealsToDiary(t *testing.T) {
	date := time.Date(2019, 5, 1, 5, 0, 3, 0, time.UTC)
	meals := meals.Meals{
		{ID: 1, Date: date, StressLevel: 2, Foods: []meals.Food{
			{Condition: meals.Cooked, Ingredient: meals.Ingredient{Name: "Ingredient 1"}},
			{Condition: meals.Raw, Ingredient: meals.Ingredient{Name: "Ingredient 2"}},
		}},
	}

	diaries := mealsToDiary(meals)

	assert.Len(t, diaries, 2)
	assert.Contains(t, diaries, Diary{Date: date, Type: DiaryFood, Content: "Ingredient 1", Severity: "cooked", Category: ""})
	assert.Contains(t, diaries, Diary{Date: date, Type: DiaryFood, Content: "Ingredient 2", Severity: "raw", Category: ""})
}

func TestSymptomsToDiary(t *testing.T) {
	date := time.Date(2019, 5, 1, 5, 0, 3, 0, time.UTC)
	events := symptoms.ConditionEvents{
		{
			Date: date,
			Conditions: []symptoms.Condition{
				{Severity: 3, Symptom: symptoms.Symptom{Name: "Symptom", SymptomCategoryID: 1}},
			},
		},
	}
	categories := symptoms.SymptomCategories{{ID: 1, Name: "Category"}}

	diaries := symptomsToDiary(events, categories)

	assert.Len(t, diaries, 1)
	assert.Contains(t, diaries, Diary{Date: date, Type: DiarySymptom, Content: "Symptom", Severity: "3", Category: "Category"})
}

func TestNotesToDiary(t *testing.T) {
	date := time.Date(2019, 5, 1, 5, 0, 3, 0, time.UTC)
	notes := []*notes.Note{{ID: 1, Date: date, Text: "Text"}}

	diaries := notesToDiary(notes)

	assert.Len(t, diaries, 1)
	assert.Contains(t, diaries, Diary{Date: date, Type: DiaryNote, Content: "Text", Severity: "", Category: ""})
}

func TestPollenToDiary(t *testing.T) {
	date := time.Date(2019, 5, 1, 5, 0, 3, 0, time.UTC)
	pollens := []pollen.PollenEvent{
		{CreatedAt: date, Pollens: []pollen.Pollen{
			{Type: pollen.Ambrosia, Intensity: pollen.HighPollen},
			{Type: pollen.Birke, Intensity: pollen.MediumToHighPollen},
		}},
	}

	diaries := pollenToDiary(pollens)

	assert.Len(t, diaries, 2)
	assert.Contains(t, diaries, Diary{Date: date, Type: DiaryPollen, Content: "", Severity: "Hohe", Category: "Ambrosia"})
	assert.Contains(t, diaries, Diary{Date: date, Type: DiaryPollen, Content: "", Severity: "Mittlere bis hohe", Category: "Birke"})
}

func TestIntakesToDiary(t *testing.T) {
	date := time.Date(2019, 5, 1, 5, 0, 3, 0, time.UTC)
	later := date.Add(time.Hour)
	intakes := []medicines.Intake{
		{Date: date, Medicine: medicines.Medicine{Name: "Medicine 1"}},
		{Date: later, Medicine: medicines.Medicine{Name: "Medicine 2"}},
	}

	diaries := intakesToDiary(intakes)

	assert.Len(t, diaries, 2)
	assert.Contains(t, diaries, Diary{Date: date, Type: DiaryIntake, Content: "Medicine 1", Severity: "", Category: ""})
	assert.Contains(t, diaries, Diary{Date: later, Type: DiaryIntake, Content: "Medicine 2", Severity: "", Category: ""})
}
