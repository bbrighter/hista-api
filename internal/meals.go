package internal

import (
	"time"

	"encore.app/entity"
)

type MealUseCase struct {
	meals       IMealsRepository
	ingredients IIngredientRepository
}

func NewMealUseCase(repo IMealsRepository, ings IIngredientRepository) MealUseCase {
	return MealUseCase{meals: repo, ingredients: ings}
}

func (uc MealUseCase) List() entity.Meals {
	return uc.meals.ListMeals()
}

func (uc MealUseCase) Get(id uint) (entity.Meal, error) {
	return uc.meals.GetMeal(id)
}

func (uc MealUseCase) Create(date *time.Time) (entity.Meal, error) {
	if date == nil {
		now := time.Now()
		date = &now
	}
	meal, err := uc.meals.CreateMeal(*date)
	return meal, err
}

func (uc MealUseCase) Delete(id uint) (ings entity.Ingredients, err error) {
	err = uc.meals.DeleteMeal(id)
	if err == nil {
		ings = uc.ingredients.ListIngredients()
	}
	return ings, err
}

func (uc MealUseCase) Patch(id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error {
	return uc.meals.PatchMeal(id, date, freshness, stressLevel, isAlone)
}
