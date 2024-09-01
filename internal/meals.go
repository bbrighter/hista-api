package internal

import (
	"time"

	"encore.app/entity"
)

type MealUseCase struct {
	Repo IMealsRepository
}

func NewMealUseCase(repo IMealsRepository) MealUseCase {
	return MealUseCase{Repo: repo}
}

func (uc MealUseCase) ListMeals() entity.Meals {
	return uc.Repo.ListMeals()
}

func (uc MealUseCase) GetMeal(id uint) (entity.Meal, error) {
	return uc.Repo.GetMeal(id)
}

func (uc MealUseCase) CreateMeal(date *time.Time) (uint, error) {
	if date == nil {
		now := time.Now()
		date = &now
	}
	id, err := uc.Repo.CreateMeal(*date)
	return id, err
}

func (uc MealUseCase) DeleteMeal(id uint) error {
	return uc.Repo.DeleteMeal(id)
}

func (uc MealUseCase) PatchMeal(id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error {
	return uc.Repo.PatchMeal(id, date, freshness, stressLevel, isAlone)
}
