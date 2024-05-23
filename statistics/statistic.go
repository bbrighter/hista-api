package statistics

import (
	"time"
)

type SymptomMeal struct {
	MealDate           time.Time `json:"foodDate"`
	ConditionEventDate time.Time `json:"symptomDate"`
	SymptomID          uint      `json:"symptomId"`
	SymptomSeverity    int       `json:"symptomSeverity"`
}

type Statistic struct {
	IngredientID  uint          `json:"ingredientId"`
	FoodCondition string        `json:"foodCondition"`
	Statistic     []SymptomMeal `json:"statistic"`
}

type Statistics struct {
	Statistics []Statistic `json:"statistics"`
}

func findFoodForSymptoms(service *Service, fromDate time.Time, toDate time.Time, symptomIds []uint) (Statistics, error) {
	type Result struct {
		MealDate           time.Time
		ConditionEventDate time.Time
		IngredientID       uint
		FoodCondition      string
		SymptomID          uint
		SymptomSeverity    int
	}

	var result []Result

	var err error = service.db.
		Select(
			"meals.date as meal_date",
			"condition_events.date as condition_event_date",
			"foods.ingredient_id",
			"foods.condition as food_condition",
			"conditions.symptom_id",
			"conditions.severity as symptom_severity",
		).
		Table("meals").
		Joins("JOIN foods ON meals.id = foods.meal_id").
		Joins("JOIN condition_events ON condition_events.date BETWEEN meals.date - interval '72 hour' AND meals.date").
		Joins("JOIN conditions ON conditions.condition_event_id = condition_events.id").
		Where("conditions.symptom_id in (?)", symptomIds).
		Where("meals.date BETWEEN ? AND ?", fromDate, toDate).
		Scan(&result).Error
	if err != nil {
		return Statistics{}, err
	}

	var statistics = []Statistic{}
	for _, res := range result {
		var exists bool = false
		var meal = SymptomMeal{MealDate: res.MealDate, ConditionEventDate: res.ConditionEventDate, SymptomID: res.SymptomID, SymptomSeverity: res.SymptomSeverity}
		for _, stat := range statistics {
			if res.IngredientID == stat.IngredientID && res.FoodCondition == stat.FoodCondition {
				exists = true
			}
			stat.Statistic = append(stat.Statistic, meal)
		}
		if !exists {
			statistics = append(statistics, Statistic{IngredientID: res.IngredientID, FoodCondition: res.FoodCondition, Statistic: []SymptomMeal{meal}})
		}

	}
	return Statistics{Statistics: statistics}, err
}
