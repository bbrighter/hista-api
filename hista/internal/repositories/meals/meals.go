package meals

import (
	"context"
	"errors"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

var ErrNoParameters = errors.New("no parameters provided")

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
	meal, err := gorm.G[entity.Meal](repo.db).
		Where("id = ?", id).
		Where("pi_id = ?", piid).
		Preload("Foods.Ingredient", nil).
		Preload("Foods", nil).
		First(ctx)
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
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	var updates = make(map[string]any)
	if date != nil {
		updates["date"] = *date
	}
	if freshness != nil {
		updates["freshness"] = *freshness
	}
	if stressLevel != nil {
		updates["stress_level"] = *stressLevel
	}
	if isAlone != nil {
		updates["is_alone"] = *isAlone
	}
	if len(updates) == 0 {
		return ErrNoParameters
	}
	tx := repo.db.Model(&entity.Meal{}).
		Where("id = ?", id).
		Where("pi_id = ?", piid).
		Updates(updates)
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return tx.Error

}

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
