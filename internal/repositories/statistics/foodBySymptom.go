package statistics

import (
	"context"
	"time"

	"encore.app/entity"
	"encore.app/generic_queries"
)

func (repo *StatisticsRepo) FindFoodForSymptoms(ctx context.Context, fromDate time.Time, toDate time.Time, symptomIds []uint) (entity.SymptomResults, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.SymptomResults{}, err
	}
	var results []entity.SymptomsResult
	subquery := repo.db.Select(
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
		Where("foods.pi_id = ?", piid).
		Where("meals.pi_id = ?", piid).
		Where("conditions.pi_id = ?", piid).
		Where("condition_events.pi_id = ?", piid).
		Group("ingredient_id, food_condition, meals.id")

	err = repo.db.
		Table("(?) as u", subquery).
		Select(
			"u.ingredient_id as ingredient_id",
			"u.food_condition as food_condition",
			"SUM(u.hours72) as hours72",
			"SUM(u.hours24) as hours24",
			"SUM(u.hours1) as hours1",
		).
		Group("ingredient_id, food_condition").
		Scan(&results).Error

	return results, err
}

func (repo *StatisticsRepo) CountFoods(ctx context.Context, relevantSymptomIds []uint) ([]entity.CountResult, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return []entity.CountResult{}, err
	}
	var countResults []entity.CountResult
	repo.db.
		Table("foods").
		Select("count(*) as count", "foods.ingredient_id as id").
		Where("ingredient_id in (?)", relevantSymptomIds).
		Where("pi_id = ?", piid).
		Group("ingredient_id").
		Scan(&countResults)
	return countResults, nil
}
