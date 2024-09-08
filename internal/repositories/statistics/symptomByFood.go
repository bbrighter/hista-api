package statistics

import (
	"time"

	"encore.app/entity"
)

func (repo *StatisticsRepo) FindSymptomsForFoods(fromDate time.Time, toDate time.Time, ingredientIds []uint) (entity.FoodResults, error) {
	var results entity.FoodResults
	subquery := repo.db.Select(
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

	var err error = repo.db.
		Table("(?) as u", subquery).
		Select(
			"u.symptom_id as symptom_id",
			"u.severity as severity",
			"SUM(u.hours72) as hours72",
			"SUM(u.hours24) as hours24",
			"SUM(u.hours1) as hours1",
		).
		Group("symptom_id, severity").
		Scan(&results).Error

	return results, err
}

func (repo *StatisticsRepo) CountSymptoms(symptomIds []uint) []entity.CountResult {
	var countResults []entity.CountResult
	repo.db.
		Table("conditions").
		Select("count(*) as count", "conditions.symptom_id as id").
		Where("symptom_id in (?)", symptomIds).
		Group("symptom_id").
		Scan(&countResults)
	return countResults
}
