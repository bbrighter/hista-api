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

func getFoods(service *Service, mealID uint) Foods {
	var foods []Food
	service.db.Where(&Food{MealID: mealID}).Preload(clause.Associations).Find(&foods)
	return foods
}

func (food *Food) createByName(service *Service, ingredientName string) (Ingredients, error) {
	var ingredients Ingredients
	if food.MealID == 0 {
		return ingredients, errors.ErrorAttributeMustBeSet("MealID")
	}
	if rows := service.db.Find(&Meal{ID: food.MealID}).RowsAffected; rows == 0 {
		return ingredients, errors.ErrorNotFound
	}
	var ingredient Ingredient
	var err error
	ingredient, err = service.createOrReplaceIngredient(ingredientName)
	if err != nil {
		return ingredients, err
	}
	food.Ingredient = ingredient
	if err := service.db.Create(food).Error; err != nil {
		return ingredients, err
	}
	service.db.Find(&ingredients)
	return ingredients, nil
}

// Create food.
// MealID and IngredientID  must be set.
func (food *Food) createByID(service *Service) (Ingredients, error) {
	var ingredients Ingredients
	if food.MealID == 0 {
		return ingredients, errors.ErrorAttributeMustBeSet("MealID")
	}
	if rows := service.db.Find(&Meal{ID: food.MealID}).RowsAffected; rows == 0 {
		return ingredients, errors.ErrorNotFound
	}
	if food.IngredientID == 0 {
		return ingredients, errors.ErrorAttributeMustBeSet("IngredientID")
	}
	if err := service.db.Create(food).Error; err != nil {
		return ingredients, err
	}
	service.db.Preload(clause.Associations).First(food)
	service.db.Find(&ingredients)
	return ingredients, nil
}

func (food *Food) delete(service *Service) error {
	if food.ID == 0 {
		return errors.ErrorIDMissing
	}
	rows := service.db.Preload("Ingredient").Find(food).RowsAffected
	if rows == 0 {
		return errors.ErrorNotFound
	}
	err := service.db.Transaction(func(tx *gorm.DB) error {
		var foodIngredient = food.Ingredient
		if err := tx.Delete(food).Error; err != nil {
			return err
		}
		return deleteIngredientIfUnused(tx, foodIngredient.ID)
	})

	return err
}

func (food *Food) changeCondition(service *Service, newCondition FoodCondition) error {
	if food.ID == 0 {
		return errors.ErrorIDMissing
	}
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
