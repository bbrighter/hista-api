package internal

import (
	"time"

	"encore.app/entity"
)

type MealUseCase struct {
	meals IMealsRepository
	ings  IIngredientRepository
}

func NewMealUseCase(repo IMealsRepository) MealUseCase {
	return MealUseCase{meals: repo}
}

func (uc MealUseCase) ListMeals() entity.Meals {
	return uc.meals.ListMeals()
}

func (uc MealUseCase) GetMeal(id uint) (entity.Meal, error) {
	return uc.meals.GetMeal(id)
}

func (uc MealUseCase) CreateMeal(date *time.Time) (uint, error) {
	if date == nil {
		now := time.Now()
		date = &now
	}
	id, err := uc.meals.CreateMeal(*date)
	return id, err
}

func (uc MealUseCase) DeleteMeal(id uint) error {
	return uc.meals.DeleteMeal(id)
}

func (uc MealUseCase) PatchMeal(id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error {
	return uc.meals.PatchMeal(id, date, freshness, stressLevel, isAlone)
}
