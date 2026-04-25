package diary

import (
	"slices"
	"strconv"
	"time"

	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/medicines"
	"encore.app/hista/internal/notes"
	"encore.app/hista/internal/pollen"
	"encore.app/hista/internal/symptoms"
)

type Diary struct {
	Date     time.Time `json:"date"`
	Type     DiaryType `json:"type"`
	Content  string    `json:"content"`
	Severity string    `json:"severity"`
	Category string    `json:"category"`
}

type DiaryType string

const (
	DiaryFood    DiaryType = "Food"
	DiarySymptom DiaryType = "Symptom"
	DiaryNote    DiaryType = "Note"
	DiaryPollen  DiaryType = "Pollen"
	DiaryIntake  DiaryType = "Intake"
)

func mealsToDiary(m meals.Meals) []Diary {
	var diaries = []Diary{}
	for _, meal := range m {
		for _, food := range meal.Foods {
			var diary = Diary{
				Date:     meal.Date,
				Type:     DiaryFood,
				Content:  food.Ingredient.Name,
				Severity: string(food.Condition),
			}
			diaries = append(diaries, diary)
		}
	}
	return diaries
}

func symptomsToDiary(e symptoms.ConditionEvents, c symptoms.SymptomCategories) []Diary {
	findSymptomCategoryById := func(id uint, categories symptoms.SymptomCategories) string {
		categoryIndex := slices.IndexFunc(categories, func(cat symptoms.SymptomCategory) bool {
			return cat.ID == id
		})
		return categories[categoryIndex].Name
	}

	var diaries = []Diary{}
	for _, event := range e {
		for _, cond := range event.Conditions {
			var diary = Diary{
				Date:     event.Date,
				Type:     DiarySymptom,
				Content:  cond.Symptom.Name,
				Severity: strconv.Itoa(int(cond.Severity)),
				Category: findSymptomCategoryById(cond.Symptom.SymptomCategoryID, c),
			}
			diaries = append(diaries, diary)
		}
	}
	return diaries
}

func notesToDiary(n []*notes.Note) []Diary {
	diaries := []Diary{}
	for _, note := range n {
		var diary = Diary{
			Date:    note.Date,
			Type:    DiaryNote,
			Content: note.Text,
		}
		diaries = append(diaries, diary)
	}
	return diaries
}

func pollenToDiary(p []pollen.PollenEvent) []Diary {
	var diaries = []Diary{}
	for _, event := range p {
		for _, pol := range event.Pollens {
			var diary = Diary{
				Date:     event.CreatedAt,
				Type:     DiaryPollen,
				Category: string(pol.Type),
				Severity: pol.Intensity.String(),
			}
			diaries = append(diaries, diary)
		}
	}
	return diaries
}

func intakesToDiary(i []medicines.Intake) []Diary {
	diaries := []Diary{}
	for _, i := range i {
		var diary = Diary{
			Date:    i.Date,
			Type:    DiaryIntake,
			Content: i.Medicine.Name,
		}
		diaries = append(diaries, diary)
	}
	return diaries
}
