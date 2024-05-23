package meals

import (
	"encore.app/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Food struct {
	ID           uint
	Ingredient   Ingredient
	IngredientID uint
	Condition    FoodCondition
	MealID       uint
}
type Foods []Food

type FoodCondition string

const (
	Raw    FoodCondition = "raw"
	Cooked FoodCondition = "cooked"
)

func (service Service) getFoods(mealID uint) Foods {
	var foods []Food
	service.db.Where(&Food{MealID: mealID}).Preload(clause.Associations).Find(&foods)
	return foods
}

func (service Service) createFood(mealID uint, condition FoodCondition, ingredientName string) (Food, error) {
	if rows := service.db.Find(&Meal{ID: mealID}).RowsAffected; rows == 0 {
		return Food{}, errors.ErrorNotFound
	}
	var ingredient Ingredient
	var err error
	var food Food
	ingredient, err = service.createOrReplaceIngredient(ingredientName)
	if err != nil {
		return food, err
	}
	food = Food{
		Ingredient: ingredient,
		Condition:  condition,
		MealID:     mealID,
	}
	err = service.db.Create(&food).Error

	return food, err
}

func (service Service) deleteFood(food Food) error {
	service.db.Preload("Ingredient").Find(&food)
	err := service.db.Transaction(func(tx *gorm.DB) error {
		var foodIngredient = food.Ingredient
		if err := tx.Delete(&food).Error; err != nil {
			return err
		}
		return deleteIngredientIfUnused(tx, foodIngredient.ID)
	})

	return err
}

func (service Service) changeFoodCondition(food Food, newCondition FoodCondition) error {
	food.Condition = newCondition
	tx := service.db.Where(&Food{ID: food.ID}).Updates(Food{Condition: newCondition})
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

func stringToFoodCondition(str string) (FoodCondition, error) {
	var err error
	var condition FoodCondition
	switch str {
	case "raw":
		condition = Raw
	case "cooked":
		condition = Cooked
	default:
		err = errors.NewError("Invalid condition: "+str, 400)
	}
	return condition, err
}
