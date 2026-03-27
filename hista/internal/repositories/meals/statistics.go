package meals

import (
	"context"
	"fmt"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
	"gorm.io/gorm"
)

var allowedTruncateUnits = map[string]bool{
	"minute":  true,
	"hour":    true,
	"day":     true,
	"week":    true,
	"month":   true,
	"quarter": true,
	"year":    true,
}

// List nutrition per truncateUnit, which is a subset of the values allowed in Postgres.
// See the map allowedTruncateUnits
func (r *MealRepository) SelectAggregatedNutrition(ctx context.Context, truncateUnit string, from *time.Time, to *time.Time) (entity.NutritionStatistics, error) {
	var nutrition entity.NutritionStatistics

	if _, ok := allowedTruncateUnits[truncateUnit]; !ok {
		return nutrition, errors.BadRequestf("invalid truncateUnit %s", truncateUnit)
	}

	dateTrunc := fmt.Sprintf("date_trunc('%s', meals.date)", truncateUnit)
	tx := r.db.Model(&entity.Meal{}).Session(&gorm.Session{SkipDefaultTransaction: true}).
		Joins("JOIN foods ON meals.id = foods.meal_id").
		Joins("JOIN ingredients ON foods.ingredient_id = ingredients.id").
		Where("ingredients.nutrition_protein IS NOT NULL").
		Where("ingredients.nutrition_carbohydrate IS NOT NULL").
		Where("ingredients.nutrition_fat IS NOT NULL").
		Where("ingredients.nutrition_fiber IS NOT NULL").
		Where("foods.amount IS NOT NULL")

	if from != nil {
		tx = tx.Where("meals.date >= ?", *from)
	}
	if to != nil {
		tx = tx.Where("meals.date <= ? ", *to)
	}

	err := tx.Select(
		dateTrunc+" as date",
		"sum(foods.amount * ingredients.nutrition_protein)/100 as protein",
		"sum(foods.amount * ingredients.nutrition_carbohydrate)/100 as carbohydrate",
		"sum(foods.amount * ingredients.nutrition_fat)/100 as fat",
		"sum(foods.amount * ingredients.nutrition_fiber)/100 as fiber",
	).Group(dateTrunc).
		Scan(&nutrition).Error

	return nutrition, err
}
