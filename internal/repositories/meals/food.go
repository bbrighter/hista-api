package meals

import (
	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *MealRepository) ListFoods(mealId uint) entity.Foods {
	var foods entity.Foods
	repo.db.Where(&entity.Food{MealID: mealId}).Preload(clause.Associations).Find(&foods)
	return foods
}

func (repo *MealRepository) CreateFoodByName(mealId uint, ingredientName string) (uint, error) {
	var food entity.Food
	if rows := repo.db.Find(&entity.Meal{ID: mealId}).RowsAffected; rows == 0 {
		return 0, errors.ErrorNotFound
	}
	var ingredient entity.Ingredient
	var err error
	ingredient, err = repo.CreateOrReplaceIngredient(ingredientName)
	if err != nil {
		return 0, err
	}
	food.MealID = mealId
	food.Ingredient = ingredient
	err = repo.db.Create(&food).Error
	return food.ID, err
}

func (repo *MealRepository) CreateFoodByID(mealId uint, ingredientId uint) (uint, error) {
	var food entity.Food
	if rows := repo.db.Find(&entity.Meal{ID: mealId}).RowsAffected; rows == 0 {
		return 0, errors.ErrorNotFound
	}
	food.MealID = mealId
	food.IngredientID = ingredientId
	if err := repo.db.Create(&food).Error; err != nil {
		return 0, err
	}
	return food.ID, nil
}

func (repo *MealRepository) DeleteFood(foodId uint) error {
	var food = entity.Food{ID: foodId}
	rows := repo.db.Preload("Ingredient").Find(&food).RowsAffected
	if rows == 0 {
		return errors.ErrorNotFound
	}
	var err error = repo.db.Transaction(func(tx *gorm.DB) error {
		var foodIngredient = food.Ingredient
		if err := tx.Delete(&food).Error; err != nil {
			return err
		}
		return deleteIngredientIfUnused(tx, foodIngredient.ID)
	})
	return err
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
	repo.db.Preload("ingredients").First(&food)
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
