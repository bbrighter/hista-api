package statistics

import (
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
				Date:    meal.Date.Truncate(time.Hour * 24),
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
				Date:    event.Date.Truncate(time.Hour * 24),
				Hour:    event.Date.Hour(),
				Type:    Symptom,
				Content: cond.Symptom.Name,
			}
			diaries = append(diaries, diary)
		}
	}
	return diaries
}

func getDiaryData(service *Service) (meals.Meals, symptoms.ConditionEvents) {
	var meals meals.Meals
	service.db.Preload("Foods.Ingredient").Preload("Foods").Find(&meals)
	var events symptoms.ConditionEvents
	service.db.Preload("Conditions.Symptom").Preload("Conditions").Find(&events)
	return meals, events

}
