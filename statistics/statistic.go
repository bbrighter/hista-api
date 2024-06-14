package statistics

import (
	"time"
)

type SymptomsResult struct {
	IngredientID  uint
	FoodCondition string
	Hours72       int
	Hours24       int
	Hours1        int
}

func countFoodBySymptoms(service *Service, fromDate time.Time, toDate time.Time, symptomIds []uint) ([]SymptomsResult, error) {
	var result []SymptomsResult
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

	var err error = service.db.
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

type StatisticsByFood struct {
	IngredientID  uint   `json:"ingredientId"`
	FoodCondition string `json:"foodCondition"`
	Hours72       int    `json:"hours72"`
	Hours24       int    `json:"hours24"`
	Hours1        int    `json:"hours1"`
	Count         int    `json:"count"`
}

type FoodStatisticsResponse struct {
	Statistics []StatisticsByFood `json:"statistics"`
}

func findFoodForSymptoms(service *Service, fromDate time.Time, toDate time.Time, symptomIds []uint) (FoodStatisticsResponse, error) {
	var err error
	var results []SymptomsResult
	results, err = countFoodBySymptoms(service, fromDate, toDate, symptomIds)
	if err != nil {
		return FoodStatisticsResponse{}, err
	}

	var stats = []StatisticsByFood{}
	for _, res := range results {
		var stat = StatisticsByFood{
			IngredientID:  res.IngredientID,
			FoodCondition: res.FoodCondition,
			Hours72:       res.Hours72,
			Hours24:       res.Hours24,
			Hours1:        res.Hours1,
		}
		stats = append(stats, stat)
	}
	return FoodStatisticsResponse{Statistics: stats}, nil
}

type FoodResult struct {
	SymptomID uint
	Severity  int
	Hours72   int
	Hours24   int
	Hours1    int
	Count     int64
}

func countSymptomsByFood(service *Service, fromDate time.Time, toDate time.Time, ingredientIds []uint) ([]FoodResult, error) {
	var result []FoodResult
	subquery := service.db.Select(
		"conditions.symptom_id as symptom_id",
		"conditions.severity as severity",
		"condition_events.id as condition_event_id",
		"MAX(CASE WHEN condition_events.date BETWEEN meals.date AND meals.date + interval '72 hour' THEN 1 ELSE 0 END) as hours72",
		"MAX(CASE WHEN condition_events.date BETWEEN meals.date AND meals.date + interval '24 hour' THEN 1 ELSE 0 END) as hours24",
		"MAX(CASE WHEN condition_events.date BETWEEN meals.date AND meals.date + interval '1 hour'  THEN 1 ELSE 0 END) as hours1",
	).
		Table("meals").
		Joins("JOIN foods ON meals.id = foods.meal_id").
		Joins("JOIN condition_events ON condition_events.date BETWEEN meals.date  AND meals.date + interval '72 hour'").
		Joins("JOIN conditions ON conditions.condition_event_id = condition_events.id").
		Where("foods.ingredient_id in (?)", ingredientIds).
		Where("condition_events.date BETWEEN ? AND ?", fromDate, toDate).
		Group("symptom_id, severity, condition_events.id")

	var err error = service.db.
		Table("(?) as u", subquery).
		Select(
			"u.symptom_id as symptom_id",
			"u.severity as severity",
			"SUM(u.hours72) as hours72",
			"SUM(u.hours24) as hours24",
			"SUM(u.hours1) as hours1",
		).
		Group("symptom_id, severity").
		Scan(&result).Error
	return result, err
}

type CountResult struct {
	ID    uint
	Count int64
}

func countSymptoms(service *Service, relevantSymptomIds []uint) []CountResult {
	var countResults []CountResult
	service.db.Debug().
		Table("conditions").
		Select("count(*) as count", "conditions.symptom_id as id").
		Where("symptom_id in (?)", relevantSymptomIds).
		Group("symptom_id").
		Scan(&countResults)
	return countResults
}

type StatisticBySymptom struct {
	SymptomID uint  `json:"symptomId"`
	Severity  int   `json:"severity"`
	Hours72   int   `json:"hours72"`
	Hours24   int   `json:"hours24"`
	Hours1    int   `json:"hours1"`
	Count     int64 `json:"count"`
}

type SymptomStatisticsResponse struct {
	Statistics []StatisticBySymptom `json:"statistics"`
}

func findSymptomsForFoods(service *Service, fromDate time.Time, toDate time.Time, ingredientIds []uint) (SymptomStatisticsResponse, error) {
	var err error
	var results []FoodResult
	results, err = countSymptomsByFood(service, fromDate, toDate, ingredientIds)
	if err != nil {
		return SymptomStatisticsResponse{}, err
	}

	var stats = []StatisticBySymptom{}
	for _, res := range results {
		var stat = StatisticBySymptom{
			SymptomID: res.SymptomID,
			Severity:  res.Severity,
			Hours72:   res.Hours72,
			Hours24:   res.Hours24,
			Hours1:    res.Hours1,
		}
		stats = append(stats, stat)
	}
	return SymptomStatisticsResponse{Statistics: stats}, nil
}
