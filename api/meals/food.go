package meals

import (
	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *MealRepository) ListFoods(mealID uint) entity.Foods {
	var foods entity.Foods
	repo.db.Where(&entity.Food{MealID: mealID}).Preload(clause.Associations).Find(&foods)
	return foods
}

func (repo *MealRepository) CreateFoodByName(mealId uint, ingredientName string) error {
	if rows := repo.db.Find(&entity.Meal{ID: mealId}).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	var ingredient entity.Ingredient
	var err error
	ingredient, err = repo.CreateOrReplaceIngredient(ingredientName)
	if err != nil {
		return err
	}
	if err := repo.db.Create(&entity.Food{MealID: mealId, Ingredient: ingredient}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MealRepository) CreateFoodByID(mealId uint, ingredientId uint) error {
	if rows := repo.db.Find(&entity.Meal{ID: mealId}).RowsAffected; rows == 0 {
		return errors.ErrorNotFound
	}
	var food = entity.Food{MealID: mealId, IngredientID: ingredientId}
	if err := repo.db.Create(&food).Error; err != nil {
		return err
	}
	repo.db.Preload(clause.Associations).First(&food) // Why is this needed?
	return nil
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

func (repo *MealRepository) ChangeCondition(foodId uint, newCondition entity.FoodCondition) error {
	tx := repo.db.Where(&entity.Food{ID: foodId}).Updates(entity.Food{Condition: newCondition})
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
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
		err = errors.NewError("Invalid condition: "+str, 400)
	}
	return condition, err
}
