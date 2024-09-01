package entity

import (
	"slices"
	"sort"
	"strconv"
	"time"

	"encore.app/pollen"
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
	FoodType    DiaryType = "Food"
	SymptomType DiaryType = "Symptom"
	NoteType    DiaryType = "Note"
	PollenType  DiaryType = "Pollen"
)

func CreateRawDiary(meals Meals, events ConditionEvents, cats SymptomCategories, notes Notes, pollens pollen.PollenEvents) []RawDiary {
	var diaries = []RawDiary{}
	for _, meal := range meals {
		for _, food := range meal.Foods {
			var diary = RawDiary{
				Date:     meal.Date,
				Type:     FoodType,
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
				Type:     SymptomType,
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
			Type:    NoteType,
			Content: note.Text,
		}
		diaries = append(diaries, diary)
	}
	for _, event := range pollens {
		for _, pol := range event.Pollens {
			var diary = RawDiary{
				Date:     event.CreatedAt,
				Type:     PollenType,
				Category: string(pol.Type),
				Severity: pol.Intensity.String(),
			}
			diaries = append(diaries, diary)
		}
	}
	sort.Slice(diaries, func(i, j int) bool {
		return diaries[j].Date.Before(diaries[i].Date)
	})
	return diaries
}

func findSymptomCategoryById(id uint, categories SymptomCategories) string {
	categoryIndex := slices.IndexFunc(categories, func(cat SymptomCategory) bool {
		return cat.ID == id
	})
	return categories[categoryIndex].Name
}
