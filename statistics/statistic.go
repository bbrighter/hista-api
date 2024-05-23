package statistics

import (
	"slices"
	"time"
)

type SymptomMeal struct {
	MealDate time.Time `json:"foodDate"`
	// ConditionEventDate time.Time `json:"symptomDate"`
	IngredientID  uint   `json:"ingredientId"`
	FoodCondition string `json:"foodCondition"`
}

type Statistic struct {
	SymptomID       uint          `json:"symptomId"`
	SymptomSeverity int           `json:"symptomSeverity"`
	SymptomMeals    []SymptomMeal `json:"statistic"`
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

	var err error = service.db.Debug().
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
	for _, id := range symptomIds {
		var stat = Statistic{
			SymptomID:    id,
			SymptomMeals: []SymptomMeal{}}
		statistics = append(statistics, stat)
	}
	for _, res := range result {
		var symptomStatistic = SymptomMeal{
			MealDate: res.MealDate,
			// ConditionEventDate: res.ConditionEventDate,
			IngredientID:  res.IngredientID,
			FoodCondition: res.FoodCondition,
		}
		var statisticIndex int = slices.IndexFunc(statistics, func(s Statistic) bool { return s.SymptomID == res.SymptomID })
		statistics[statisticIndex].SymptomMeals = append(statistics[statisticIndex].SymptomMeals, symptomStatistic)
		statistics[statisticIndex].SymptomSeverity = res.SymptomSeverity
	}
	return Statistics{Statistics: statistics}, err
}
