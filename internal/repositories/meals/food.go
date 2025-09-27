package meals

import (
	"context"

	"encore.app/entity"
	"encore.app/errors"
	"encore.app/generic_queries"
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
	return gorm.G[entity.Food](repo.db).Create(ctx, food)
}

func (repo *MealRepository) CreateFoodByID(ctx context.Context, food *entity.Food) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	food.PIID = piid
	if food.MealID == 0 || food.IngredientID == 0 || food.Condition == "" {
		return errors.ErrorAttributeMustBeSet("MealID or IngredientID or Condition")
	}
	if rows := repo.db.Find(&entity.Meal{ID: food.MealID, PIID: piid}).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	return repo.db.Create(&food).Error
}

func (repo *MealRepository) DeleteFood(ctx context.Context, foodId uint) error {
	// TODO: Delete ingredients as well if last!
	return generic_queries.Delete[*entity.Food](ctx, repo.db, foodId)
}

func (repo *MealRepository) ChangeCondition(ctx context.Context, foodId uint, condition entity.FoodCondition) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	tx := repo.db.Where(&entity.Food{ID: foodId, PIID: piid}).Updates(entity.Food{Condition: condition})
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
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

func StringToFoodCondition(str string) (entity.FoodCondition, error) {
	var err error
	var condition entity.FoodCondition
	switch str {
	case "raw":
		condition = entity.Raw
	case "cooked":
		condition = entity.Cooked
	default:
		err = errors.BadRequest("invalid condition")
	}
	return condition, err
}
