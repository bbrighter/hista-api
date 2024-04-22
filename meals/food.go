package meals

import (
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

type FoodCondition string

const (
	Raw    FoodCondition = "raw"
	Cooked FoodCondition = "cooked"
)

func (service Service) getFoods(mealID uint) []Food {
	var foods []Food
	service.db.Where(&Food{MealID: mealID}).Preload(clause.Associations).Find(&foods)
	return foods
}

func (service Service) createFood(mealID uint, condition FoodCondition, ingredientName string) (Food, error) {
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
	err := service.db.Transaction(func(tx *gorm.DB) error {
		var foodIngredient = food.Ingredient
		if err := tx.Delete(&food).Error; err != nil {
			return err
		}
		return deleteIngredientIfUnused(tx, foodIngredient.ID)
	})

	return err
}
