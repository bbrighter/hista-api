package statistics

import (
	"slices"
	"sort"
	"strconv"
	"time"

	"encore.app/meals"
	"encore.app/symptoms"
)

type RawDiary struct {
	Date     time.Time `json:"date"`
	Hour     int       `json:"hour"`
	Type     DiaryType `json:"type"`
	Content  string    `json:"content"`
	Severity string    `json:"severity"`
	Category string    `json:"category"`
}

type DiaryType string

const (
	Food    DiaryType = "Food"
	Symptom DiaryType = "Symptom"
)

func diaryFrom(meals meals.Meals, events symptoms.ConditionEvents, symptomCategories symptoms.SymptomCategories) []RawDiary {
	var diaries = []RawDiary{}
	for _, meal := range meals {
		for _, food := range meal.Foods {
			var diary = RawDiary{
				Date:     meal.Date,
				Hour:     meal.Date.Hour(),
				Type:     Food,
				Content:  food.Ingredient.Name,
				Severity: string(food.Condition),
				Category: "",
			}
			diaries = append(diaries, diary)
		}

	}
	for _, event := range events {
		for _, cond := range event.Conditions {
			var categoryId uint = cond.Symptom.SymptomCategoryID
			categoryIndex := slices.IndexFunc(symptomCategories, func(cat symptoms.SymptomCategory) bool {
				return cat.ID == categoryId
			})
			var diary = RawDiary{
				Date:     event.Date,
				Hour:     event.Date.Hour(),
				Type:     Symptom,
				Content:  cond.Symptom.Name,
				Severity: strconv.Itoa(int(cond.Severity)),
				Category: symptomCategories[categoryIndex].Name,
			}
			diaries = append(diaries, diary)
		}
	}
	sort.Slice(diaries, func(i, j int) bool {
		return diaries[j].Date.Before(diaries[i].Date)
	})
	return diaries
}
