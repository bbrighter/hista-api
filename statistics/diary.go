package statistics

import (
	"slices"
	"sort"
	"strconv"
	"time"

	"encore.app/entity"
	"encore.app/pollen"
)

type Input struct {
	Meals      entity.Meals
	Events     entity.ConditionEvents
	Categories entity.SymptomCategories
	Notes      entity.Notes
	Pollens    []pollen.PollenEvent
}

type RawDiary struct {
	Date     time.Time `json:"date"`
	Type     DiaryType `json:"type"`
	Content  string    `json:"content"`
	Severity string    `json:"severity"`
	Category string    `json:"category"`
}

type DiaryType string

const (
	Food    DiaryType = "Food"
	Symptom DiaryType = "Symptom"
	Note    DiaryType = "Note"
	Pollen  DiaryType = "Pollen"
)

func createRawDiary(input Input) []RawDiary {
	var diaries = []RawDiary{}
	for _, meal := range input.Meals {
		for _, food := range meal.Foods {
			var diary = RawDiary{
				Date:     meal.Date,
				Type:     Food,
				Content:  food.Ingredient.Name,
				Severity: string(food.Condition),
			}
			diaries = append(diaries, diary)
		}
	}
	for _, event := range input.Events {
		for _, cond := range event.Conditions {
			var diary = RawDiary{
				Date:     event.Date,
				Type:     Symptom,
				Content:  cond.Symptom.Name,
				Severity: strconv.Itoa(int(cond.Severity)),
				Category: findSymptomCategoryById(cond.Symptom.SymptomCategoryID, input.Categories),
			}
			diaries = append(diaries, diary)
		}
	}
	for _, note := range input.Notes {
		var diary = RawDiary{
			Date:    note.Date,
			Type:    Note,
			Content: note.Text,
		}
		diaries = append(diaries, diary)
	}
	for _, event := range input.Pollens {
		for _, pol := range event.Pollens {
			var diary = RawDiary{
				Date:     event.CreatedAt,
				Type:     Pollen,
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

func findSymptomCategoryById(id uint, categories entity.SymptomCategories) string {
	categoryIndex := slices.IndexFunc(categories, func(cat entity.SymptomCategory) bool {
		return cat.ID == id
	})
	return categories[categoryIndex].Name
}
