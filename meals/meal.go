package meals

import (
	"errors"
	"time"

	"encore.dev/beta/errs"
	"gorm.io/gorm"
)

type Meal struct {
	ID    uint
	Date  time.Time
	Foods []Food `gorm:"constraint:OnDelete:CASCADE"`
}

func (service Service) createMeal(date time.Time) (uint, error) {
	var meal = Meal{Date: date}
	service.db.Create(&meal)
	return meal.ID, nil
}

func (service Service) getMeals() []Meal {
	var meals []Meal
	service.db.Find(&meals)
	return meals
}

func (service Service) getMeal(id uint) (Meal, error) {
	var meal = Meal{ID: id}
	if service.db.Preload("Foods.Ingredient").Preload("Foods").Find(&meal).RowsAffected == 0 {
		return Meal{}, &errs.Error{Code: errs.NotFound}
	}
	return meal, nil
}

func (service Service) deleteMeal(mealId uint) error {
	var meal Meal
	if service.db.Find(&meal).RowsAffected == 0 {
		return errors.New("not found")
	}
	return service.db.Transaction(func(tx *gorm.DB) error {
		var foods []Food
		tx.Where(&Food{MealID: mealId}).Find(&foods)
		if err := tx.Delete(&Food{}, Food{MealID: mealId}).Error; err != nil {
			return err
		}
		for _, food := range foods {
			deleteIngredientIfUnused(tx, food.IngredientID)
		}
		return tx.Delete(&Meal{ID: mealId}).Error
	})
}
