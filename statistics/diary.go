package statistics

import (
	"slices"
	"sort"
	"strconv"
	"time"

	"encore.app/meals"
	"encore.app/symptoms"
)

type Input struct {
	Meals      meals.Meals
	Events     symptoms.ConditionEvents
	Categories symptoms.SymptomCategories
}

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

func diaryFrom(input Input) []RawDiary {
	var diaries = []RawDiary{}
	for _, meal := range input.Meals {
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
	for _, event := range input.Events {
		for _, cond := range event.Conditions {
			var diary = RawDiary{
				Date:     event.Date,
				Hour:     event.Date.Hour(),
				Type:     Symptom,
				Content:  cond.Symptom.Name,
				Severity: strconv.Itoa(int(cond.Severity)),
				Category: findSymptomCategoryById(cond.Symptom.SymptomCategoryID, input.Categories),
			}
			diaries = append(diaries, diary)
		}
	}
	sort.Slice(diaries, func(i, j int) bool {
		return diaries[j].Date.Before(diaries[i].Date)
	})
	return diaries
}

type DiarySymptom struct {
	Name     string
	Severity int
	Category string
}

type MealBasedDiary struct {
	Date              time.Time
	Hour              int
	Food              string
	Condition         string
	SymptomsWithin1h  []DiarySymptom
	SymptomsWithin12h []DiarySymptom
	SymptomsWithin24h []DiarySymptom
}

func mealBasedDiary(input Input) []MealBasedDiary {
	var diaries []MealBasedDiary
	for _, meal := range input.Meals {
		for _, event := range input.Events {
			var symps1 []DiarySymptom = input.findSymptomsWithin(time.Hour, event, meal)
			var symps12 []DiarySymptom = input.findSymptomsWithin(12*time.Hour, event, meal)
			var symps24 []DiarySymptom = input.findSymptomsWithin(24*time.Hour, event, meal)
			for _, food := range meal.Foods {
				var diary = MealBasedDiary{
					Date:              meal.Date,
					Hour:              meal.Date.Hour(),
					Food:              food.Ingredient.Name,
					Condition:         string(food.Condition),
					SymptomsWithin1h:  symps1,
					SymptomsWithin12h: symps12,
					SymptomsWithin24h: symps24,
				}
				diaries = append(diaries, diary)
			}
		}
	}
	return diaries
}

func (input Input) findSymptomsWithin(
	duration time.Duration,
	event symptoms.ConditionEvent,
	meal meals.Meal,
) []DiarySymptom {
	var symps []DiarySymptom
	for _, con := range event.Conditions {
		if meal.Date.Sub(event.Date) < duration {
			var symptom = DiarySymptom{
				Name:     con.Symptom.Name,
				Severity: int(con.Severity),
				Category: findSymptomCategoryById(con.Symptom.SymptomCategoryID, input.Categories),
			}
			symps = append(symps, symptom)
		}
	}
	return symps
}

func findSymptomCategoryById(id uint, categories symptoms.SymptomCategories) string {
	categoryIndex := slices.IndexFunc(categories, func(cat symptoms.SymptomCategory) bool {
		return cat.ID == id
	})
	return categories[categoryIndex].Name
}
