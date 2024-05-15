package statistics

import (
	"sort"
	"time"

	"encore.app/meals"
	"encore.app/symptoms"
)

type RawDiary struct {
	Date    time.Time `json:"date"`
	Hour    int       `json:"hour"`
	Type    Category  `json:"type"`
	Content string    `json:"content"`
}

type Category string

const (
	Food    Category = "Food"
	Symptom Category = "Symptom"
)

func diaryFrom(meals meals.Meals, events symptoms.ConditionEvents) []RawDiary {
	var diaries = []RawDiary{}
	for _, meal := range meals {
		for _, food := range meal.Foods {
			var diary = RawDiary{
				Date:    meal.Date,
				Hour:    meal.Date.Hour(),
				Type:    Food,
				Content: food.Ingredient.Name,
			}
			diaries = append(diaries, diary)
		}

	}
	for _, event := range events {
		for _, cond := range event.Conditions {
			var diary = RawDiary{
				Date:    event.Date,
				Hour:    event.Date.Hour(),
				Type:    Symptom,
				Content: cond.Symptom.Name,
			}
			diaries = append(diaries, diary)
		}
	}
	sort.Slice(diaries, func(i, j int) bool {
		return diaries[j].Date.Before(diaries[i].Date)
	})
	return diaries
}
