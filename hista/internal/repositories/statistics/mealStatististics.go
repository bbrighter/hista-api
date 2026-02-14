package statistics

import (
	"context"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

func mealDateBetween(fromDate, toDate time.Time) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("meals.date between ? and ?", fromDate, toDate)
	}
}

func ingredientIdIs(id uint) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("foods.ingredient_id = ?", id)
	}
}

func (repo *StatisticsRepo) SymptomsAfterIngredients(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (entity.FoodResults, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.FoodResults{}, err
	}
	selectFields := []string{
		"conditions.symptom_id as symptom_id",
		"conditions.severity as severity",
		"meals.id as meal_id",
	}
	selectFields = append(selectFields, buildTimewindowSelects()...)
	var foodResults []entity.FoodResult
	subquery := repo.db.
		Model(entity.Meal{}).
		Scopes(mealDateBetween(fromDate, toDate), ingredientIdIs(ingredientId)).
		Where("meals.pi_id = ?", piid).
		Joins("JOIN foods on meals.id = foods.meal_id and meals.pi_id = foods.pi_id").
		Joins("JOIN condition_events ON condition_events.date BETWEEN meals.date  AND meals.date + interval '72 hour'").
		Where("condition_events.pi_id = ?", piid).
		Joins("JOIN conditions ON conditions.condition_event_id = condition_events.id and condition_events.pi_id = conditions.pi_id").
		Select(selectFields).
		Group("meals.id").Group("conditions.symptom_id").Group("conditions.severity")

	grouped_select_fields := []string{
		"u.symptom_id as symptom_id",
		"u.severity as severity"}
	grouped_select_fields = append(grouped_select_fields, buildTimeWindowSum()...)
	err = repo.db.
		Table("(?) as u", subquery).
		Select(grouped_select_fields).
		Group("symptom_id").Group("severity").
		Scan(&foodResults).Error
	return foodResults, err
}

func (r *StatisticsRepo) CountMealsWithIngredients(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (int64, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	var counts int64
	err = r.db.
		Model(entity.Meal{}).
		Joins("JOIN foods ON foods.meal_id = meals.id and foods.pi_id = meals.pi_id").
		Scopes(mealDateBetween(fromDate, toDate), ingredientIdIs(ingredientId)).
		Where("meals.pi_id = ?", piid).
		Select("count(distinct meals.id) as count").
		Count(&counts).Error

	return counts, err
}
