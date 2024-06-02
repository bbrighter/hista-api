package statistics

import (
	"slices"
	"sort"
	"strconv"
	"time"

	"encore.app/meals"
	"encore.app/notes"
	"encore.app/pollen"
	"encore.app/symptoms"
)

type Input struct {
	Meals      meals.Meals
	Events     symptoms.ConditionEvents
	Categories symptoms.SymptomCategories
	Notes      notes.Notes
	Pollens    pollen.Pollens
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
				Category: "",
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
			Date:     note.Date,
			Type:     Note,
			Content:  note.Text,
			Severity: "",
			Category: "",
		}
		diaries = append(diaries, diary)
	}
	for _, pol := range input.Pollens {
		var relevantPollen = make(map[string]pollen.PollenLoad)
		relevantPollen["Ambrosia"] = pol.Ambrosia
		relevantPollen["Beifuss"] = pol.Beifuss
		relevantPollen["Birke"] = pol.Birke
		relevantPollen["Erle"] = pol.Erle
		relevantPollen["Esche"] = pol.Esche
		relevantPollen["Graeser"] = pol.Graeser
		relevantPollen["Hasel"] = pol.Hasel
		relevantPollen["Roggen"] = pol.Roggen

		for key, value := range relevantPollen {
			if value > pollen.No {
				var diary = RawDiary{
					Date:     pol.CreatedAt,
					Type:     Pollen,
					Category: key,
					Severity: value.String(),
				}
				diaries = append(diaries, diary)
			}
		}
	}
	sort.Slice(diaries, func(i, j int) bool {
		return diaries[j].Date.Before(diaries[i].Date)
	})
	return diaries
}

func findSymptomCategoryById(id uint, categories symptoms.SymptomCategories) string {
	categoryIndex := slices.IndexFunc(categories, func(cat symptoms.SymptomCategory) bool {
		return cat.ID == id
	})
	return categories[categoryIndex].Name
}
