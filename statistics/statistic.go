package statistics

import (
	"time"
)

type StatisticByFood struct {
	IngredientID  uint   `json:"ingredientId"`
	FoodCondition string `json:"foodCondition"`
	Hours72       int    `json:"hours72"`
	Hours24       int    `json:"hours24"`
	Hours1        int    `json:"hours1"`
}

type StatisticsResponse struct {
	Statistics []StatisticByFood `json:"statistics"`
}

type Result struct {
	IngredientID  uint
	FoodCondition string
	Hours72       int
	Hours24       int
	Hours1        int
}

func getFoodsAndSymptoms(service *Service, fromDate time.Time, toDate time.Time, symptomIds []uint) ([]Result, error) {
	var result []Result
	subquery := service.db.Select(
		"foods.ingredient_id as ingredient_id",
		"foods.condition as food_condition",
		"meals.id as meal_id",
		"MAX(CASE WHEN condition_events.date BETWEEN meals.date  AND meals.date + interval '72 hour' THEN 1 ELSE 0 END) as hours72",
		"MAX(CASE WHEN condition_events.date BETWEEN meals.date  AND meals.date + interval '24 hour' THEN 1 ELSE 0 END) as hours24",
		"MAX(CASE WHEN condition_events.date BETWEEN meals.date  AND meals.date + interval '1 hour' THEN 1 ELSE 0 END) as hours1",
	).
		Table("meals").
		Joins("JOIN foods ON meals.id = foods.meal_id").
		Joins("JOIN condition_events ON condition_events.date BETWEEN meals.date  AND meals.date + interval '72 hour'").
		Joins("JOIN conditions ON conditions.condition_event_id = condition_events.id").
		Where("conditions.symptom_id in (?)", symptomIds).
		Where("meals.date BETWEEN ? AND ?", fromDate, toDate).
		Group("ingredient_id, food_condition, meals.id")

	var err error = service.db.Debug().
		Table("(?) as u", subquery).
		Select(
			"u.ingredient_id as ingredient_id",
			"u.food_condition as food_condition",
			"SUM(u.hours72) as hours72",
			"SUM(u.hours24) as hours24",
			"SUM(u.hours1) as hours1",
		).
		Group("ingredient_id, food_condition").
		Scan(&result).Error

	return result, err
}

func findFoodForSymptoms(service *Service, fromDate time.Time, toDate time.Time, symptomIds []uint) (StatisticsResponse, error) {
	var err error
	var results []Result
	results, err = getFoodsAndSymptoms(service, fromDate, toDate, symptomIds)
	if err != nil {
		return StatisticsResponse{}, err
	}

	var stats = []StatisticByFood{}
	for _, res := range results {
		println("res", res.FoodCondition, res.IngredientID, res.Hours1, res.Hours24, res.Hours72)
		var stat = StatisticByFood{
			IngredientID:  res.IngredientID,
			FoodCondition: res.FoodCondition,
			Hours72:       res.Hours72,
			Hours24:       res.Hours24,
			Hours1:        res.Hours1,
		}
		stats = append(stats, stat)
	}
	return StatisticsResponse{Statistics: stats}, nil
}
