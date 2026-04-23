package statistics

import (
	"context"
	"fmt"
	"time"

	"encore.app/hista/internal/meals"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type statisticsRepo struct {
	db *gorm.DB
}

func newStatisticsRepo(db *gorm.DB) *statisticsRepo {
	return &statisticsRepo{db: db}
}

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

func (repo *statisticsRepo) SymptomsAfterIngredients(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (FoodResults, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return FoodResults{}, err
	}
	selectFields := []string{
		"conditions.symptom_id as symptom_id",
		"conditions.severity as severity",
		"meals.id as meal_id",
	}
	selectFields = append(selectFields, buildTimeWindowSelects()...)
	var foodResults []FoodResult
	subquery := repo.db.
		Model(meals.Meal{}).
		Scopes(mealDateBetween(fromDate, toDate), ingredientIdIs(ingredientId)).
		Where("meals.pi_id = ?", piid).
		Joins("JOIN foods on meals.id = foods.meal_id").
		Joins("JOIN condition_events ON condition_events.date BETWEEN meals.date  AND meals.date + interval '72 hour'").
		Where("condition_events.pi_id = ?", piid).
		Joins("JOIN conditions ON conditions.condition_event_id = condition_events.id").
		Select(selectFields).
		Group("meals.id").Group("conditions.symptom_id").Group("conditions.severity")

	grouped_select_fields := []string{
		"u.symptom_id as symptom_id",
		"u.severity as severity"}
	grouped_select_fields = append(grouped_select_fields, buildTimeWindowSum()...)
	err = repo.db.Session(&gorm.Session{SkipDefaultTransaction: true}).
		Table("(?) as u", subquery).
		Select(grouped_select_fields).
		Group("symptom_id").Group("severity").
		Scan(&foodResults).Error
	return foodResults, err
}

func (r *statisticsRepo) CountMealsWithIngredients(ctx context.Context, fromDate time.Time, toDate time.Time, ingredientId uint) (int64, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	var counts int64
	err = r.db.
		Model(meals.Meal{}).Session(&gorm.Session{SkipDefaultTransaction: true}).
		Joins("JOIN foods ON foods.meal_id = meals.id").
		Scopes(mealDateBetween(fromDate, toDate), ingredientIdIs(ingredientId)).
		Where("meals.pi_id = ?", piid).
		Select("count(distinct meals.id) as count").
		Count(&counts).Error

	return counts, err
}

type timeWindow struct {
	name      string
	fromHours int
	toHours   int
}

func newTimeWindows() []timeWindow {
	return []timeWindow{
		{name: "hours_72", fromHours: 24, toHours: 72},
		{name: "hours_24", fromHours: 1, toHours: 24},
		{name: "hours_1", fromHours: 0, toHours: 1},
	}
}

func buildTimeWindowSelects() []string {
	windows := newTimeWindows()
	var selects []string
	for _, w := range windows {
		s := fmt.Sprintf(`MAX(
	CASE WHEN condition_events.date BETWEEN meals.date + interval '%d hour' AND meals.date + interval '%d hour' THEN 1 
	ELSE 0 END
	) as %s`, w.fromHours, w.toHours, w.name)
		selects = append(selects, s)
	}
	return selects
}

func buildTimeWindowSum() []string {
	windows := newTimeWindows()
	var selects []string
	for _, w := range windows {
		s := fmt.Sprintf("sum(u.%s) as %s", w.name, w.name)
		selects = append(selects, s)
	}
	return selects
}
