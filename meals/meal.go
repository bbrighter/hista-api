package meals

import (
	"errors"
	"time"

	"encore.dev/beta/errs"
	"gorm.io/gorm/clause"
)

type Meal struct {
	ID    uint
	Date  time.Time
	Foods []Food `gorm:"constraint:OnDelete:CASCADE"`
}

func (service Service) createMeal(foods []Food, date time.Time) (uint, error) {
	var ingredientIDs []uint
	for _, f := range foods {
		ingredientIDs = append(ingredientIDs, f.IngredientID)
	}
	var foundIngredients []Ingredient
	service.db.Find(&foundIngredients, ingredientIDs)
	for _, i := range ingredientIDs {
		var found bool = false
		for _, j := range foundIngredients {
			if i == j.ID {
				found = true
			}
		}
		if !found {
			return 0, errors.New("Used ingredient missing")
		}
	}

	var meal = Meal{
		Date:  date,
		Foods: foods,
	}
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
	if service.db.Preload(clause.Associations).Find(&meal).RowsAffected == 0 {
		return Meal{}, &errs.Error{Code: errs.NotFound}
	}
	return meal, nil
}
