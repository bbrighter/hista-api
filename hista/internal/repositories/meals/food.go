package meals

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *MealRepository) ListFoods(ctx context.Context, mealId uint) ([]*entity.Food, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return []*entity.Food{}, err
	}
	return gorm.G[*entity.Food](repo.db).Where(&entity.Food{MealID: mealId, PIID: piid}).Preload("Ingredient", nil).Find(ctx)
}

func (repo *MealRepository) CreateFoodByName(ctx context.Context, food *entity.Food, ingredientName string) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}

	if food.MealID == 0 {
		return errors.ErrorAttributeMustBeSet("MealID")
	}
	count, err := gorm.G[entity.Meal](repo.db).Where("pi_id = ?", piid).Where("id = ?", food.MealID).Count(ctx, "*")
	if err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	ingredient, err := repo.CreateOrReplaceIngredient(ctx, ingredientName)
	if err != nil {
		return err
	}
	food.Ingredient = ingredient
	food.PIID = piid
	food.IngredientPIID = piid
	food.MealPIID = piid
	return gorm.G[entity.Food](repo.db).Create(ctx, food)
}

func (repo *MealRepository) CreateFoodByID(ctx context.Context, food *entity.Food) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	food.PIID = piid
	food.IngredientPIID = piid
	food.MealPIID = piid
	if food.MealID == 0 || food.IngredientID == 0 || food.Condition == "" {
		return errors.ErrorAttributeMustBeSet("MealID or IngredientID or Condition")
	}
	rows, err := gorm.G[entity.Meal](repo.db).Where("pi_id = ?", piid).Where("id = ?", food.MealID).Count(ctx, "*")
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return gorm.G[entity.Food](repo.db).Create(ctx, food)
}

func (repo *MealRepository) DeleteFood(ctx context.Context, foodId uint) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	food, err := generic_queries.First[*entity.Food](ctx, repo.db, foodId)
	if err != nil {
		return err
	}
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := generic_queries.Delete[*entity.Food](ctx, tx, foodId); err != nil {
			return err
		}
		numberOfIngredientUsage, err := gorm.G[entity.Food](tx).
			Where("pi_id = ?", piid).
			Where("ingredient_id = ?", food.IngredientID).
			Count(ctx, "*")
		if err != nil {
			return err
		}
		if numberOfIngredientUsage == 0 {
			return generic_queries.Delete[*entity.Ingredient](ctx, tx, food.IngredientID)
		}

		return nil
	})
}

func (repo *MealRepository) UpdateFood(ctx context.Context, foodId uint, column string, value any) error {
	return generic_queries.UpdateColumn[*entity.Food](ctx, repo.db, foodId, column, value)
}

func (repo *MealRepository) GetFood(ctx context.Context, foodId uint) (entity.Food, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.Food{}, err
	}
	var food = entity.Food{ID: foodId, PIID: piid}
	repo.db.Preload(clause.Associations).First(&food)
	return food, nil
}
