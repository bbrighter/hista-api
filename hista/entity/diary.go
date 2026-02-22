package entity

import (
	"slices"
	"sort"
	"strconv"
	"time"
)

type RawDiary struct {
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

func CreateRawDiary(meals Meals, events ConditionEvents, cats SymptomCategories, notes Notes, pollens PollenEvents, intakes []*Intake) []RawDiary {
	var diaries = []RawDiary{}
	for _, meal := range meals {
		for _, food := range meal.Foods {
			var diary = RawDiary{
				Date:     meal.Date,
				Type:     DiaryFood,
				Content:  food.Ingredient.Name,
				Severity: string(food.Condition),
			}
			diaries = append(diaries, diary)
		}
	}
	for _, event := range events {
		for _, cond := range event.Conditions {
			var diary = RawDiary{
				Date:     event.Date,
				Type:     DiarySymptom,
				Content:  cond.Symptom.Name,
				Severity: strconv.Itoa(int(cond.Severity)),
				Category: findSymptomCategoryById(cond.Symptom.SymptomCategoryID, cats),
			}
			diaries = append(diaries, diary)
		}
	}
	for _, note := range notes {
		var diary = RawDiary{
			Date:    note.Date,
			Type:    DiaryNote,
			Content: note.Text,
		}
		diaries = append(diaries, diary)
	}
	for _, event := range pollens {
		for _, pol := range event.Pollens {
			var diary = RawDiary{
				Date:     event.CreatedAt,
				Type:     DiaryPollen,
				Category: string(pol.Type),
				Severity: pol.Intensity.String(),
			}
			diaries = append(diaries, diary)
		}
	}
	for _, i := range intakes {
		var diary = RawDiary{
			Date:    i.Date,
			Type:    DiaryIntake,
			Content: i.Medicine.Name,
		}
		diaries = append(diaries, diary)
	}
	sort.Slice(diaries, func(i, j int) bool {
		return diaries[j].Date.Before(diaries[i].Date)
	})
	return diaries
}

func findSymptomCategoryById(id uint, categories SymptomCategories) string {
	categoryIndex := slices.IndexFunc(categories, func(cat *SymptomCategory) bool {
		return cat.ID == id
	})
	return categories[categoryIndex].Name
}
