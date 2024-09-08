package meals

import (
	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm/clause"
)

func (repo *MealRepository) ListFoods(mealId uint) entity.Foods {
	var foods entity.Foods
	repo.db.Where(&entity.Food{MealID: mealId}).Preload(clause.Associations).Find(&foods)
	return foods
}

func (repo *MealRepository) CreateFoodByName(food *entity.Food, ingredientName string) error {
	if food.MealID == 0 {
		return errors.ErrorAttributeMustBeSet("MealID")
	}

	if rows := repo.db.Find(&entity.Meal{ID: food.MealID}).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	var ingredient entity.Ingredient
	var err error
	ingredient, err = repo.CreateOrReplaceIngredient(ingredientName)
	if err != nil {
		return err
	}
	food.Ingredient = ingredient
	err = repo.db.Create(&food).Error
	return err
}

func (repo *MealRepository) CreateFoodByID(food *entity.Food) error {
	if food.MealID == 0 || food.IngredientID == 0 || food.Condition == "" {
		return errors.ErrorAttributeMustBeSet("MealID or IngredientID or Condition")
	}
	if rows := repo.db.Find(&entity.Meal{ID: food.MealID}).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	return repo.db.Create(&food).Error
}

func (repo *MealRepository) DeleteFood(foodId uint) error {
	var food = entity.Food{ID: foodId}
	tx := repo.db.Delete(&food)
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

func (repo *MealRepository) ChangeCondition(foodId uint, condition entity.FoodCondition) error {
	tx := repo.db.Where(&entity.Food{ID: foodId}).Updates(entity.Food{Condition: condition})
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

func (repo *MealRepository) GetFood(foodId uint) entity.Food {
	var food = entity.Food{ID: foodId}
	repo.db.Preload(clause.Associations).First(&food)
	return food
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
