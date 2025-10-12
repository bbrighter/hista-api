package meals

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

func (repo *MealRepository) ListMeals(ctx context.Context) ([]*entity.Meal, error) {
	return generic_queries.List[*entity.Meal](ctx, repo.db)
}

func (repo *MealRepository) CreateMeal(ctx context.Context, meal *entity.Meal) error {
	return generic_queries.Create(ctx, repo.db, meal)
}

func (repo *MealRepository) GetMeal(ctx context.Context, id uint) (entity.Meal, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.Meal{}, err
	}
	meal, err := gorm.G[entity.Meal](repo.db).Where("pi_id = ?", piid).Preload("Foods.Ingredient", nil).Preload("Foods", nil).First(ctx)
	return meal, err
}

func (repo *MealRepository) DeleteMeal(ctx context.Context, id uint) error {
	meal, err := repo.GetMeal(ctx, id)
	if err != nil {
		return err
	}
	repo.db.Transaction(func(tx *gorm.DB) error {
		var foodIds []uint
		for _, f := range meal.Foods {
			foodIds = append(foodIds, f.ID)
		}
		_, err := gorm.G[entity.Food](tx).Where("id IN ?", foodIds).Delete(ctx)
		if err != nil {
			return err
		}

		unusedIngredients, _ := gorm.G[entity.Ingredient](tx).
			Where("pi_id = ?", meal.PIID).
			Where("NOT EXISTS (SELECT 1 FROM foods WHERE foods.ingredient_id = ingredients.id)").
			Find(ctx)

		var ingredientIds []uint
		for _, i := range unusedIngredients {
			ingredientIds = append(ingredientIds, i.ID)
		}
		_, err = gorm.G[entity.Ingredient](tx).Where("id IN ? ", ingredientIds).Delete(ctx)

		return err
	})

	return generic_queries.Delete[*entity.Meal](ctx, repo.db, id)
}

// PatchMeal a meal with parameters. Only given parameters are patched.
func (repo *MealRepository) PatchMeal(
	ctx context.Context,
	id uint,
	date *time.Time,
	freshness *entity.Freshness,
	stressLevel *uint8,
	isAlone *bool,
) error {
	var meal entity.Meal
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	if rows := repo.db.Where("pi_id = ?", piid).First(&meal, &entity.Meal{ID: id}).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	tx := repo.db.Model(&entity.Meal{ID: meal.ID})
	var updates = make(map[string]interface{})
	if date != nil {
		meal.Date = *date
		updates["date"] = *date
	}
	if freshness != nil {
		meal.Freshness = *freshness
		updates["freshness"] = *freshness
	}
	if stressLevel != nil {
		meal.StressLevel = *stressLevel
		updates["stress_level"] = *stressLevel
	}
	if isAlone != nil {
		meal.IsAlone = *isAlone
		updates["is_alone"] = *isAlone
	}
	return tx.Updates(updates).Error
}

// func GetMealsAndDependencies(db *gorm.DB) (entity.Meals, error) {
// 	var meals entity.Meals
// 	var err error = db.Preload("Foods.Ingredient").
// 		Preload("Foods").
// 		Find(&meals).Error
// 	return meals, err
// }

func (repo *MealRepository) ListMealsWithDependencies(ctx context.Context) (entity.Meals, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.Meals{}, err
	}
	var meals entity.Meals
	tx := repo.db.Preload("Foods.Ingredient").
		Where("pi_id = ?", piid).
		Preload("Foods").
		Find(&meals)
	return meals, tx.Error
}
