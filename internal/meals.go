package internal

import (
	"context"
	"time"

	"encore.app/entity"
)

type (
	IMealsRepository interface {
		ListMeals(ctx context.Context) ([]*entity.Meal, error)
		CreateMeal(ctx context.Context, meal *entity.Meal) error
		GetMeal(ctx context.Context, id uint) (entity.Meal, error)
		DeleteMeal(ctx context.Context, id uint) error
		PatchMeal(ctx context.Context, id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error
		ListMealsWithDependencies(ctx context.Context) (entity.Meals, error)
	}

	IMealUseCase interface {
		List(ctx context.Context) (entity.Meals, error)
		Create(ctx context.Context, date *time.Time) (entity.Meal, error)
		Get(ctx context.Context, id uint) (entity.Meal, error)
		Delete(ctx context.Context, id uint) (entity.Ingredients, error)
		Patch(ctx context.Context, id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error
	}
)

type MealUseCase struct {
	meals       IMealsRepository
	ingredients IIngredientRepository
}

func NewMealUseCase(repo IMealsRepository, ings IIngredientRepository) MealUseCase {
	return MealUseCase{meals: repo, ingredients: ings}
}

func (uc MealUseCase) List(ctx context.Context) (entity.Meals, error) {
	meals, err := uc.meals.ListMeals(ctx)
	return meals, errorMapper(err)
}

func (uc MealUseCase) Get(ctx context.Context, id uint) (entity.Meal, error) {
	meal, err := uc.meals.GetMeal(ctx, id)
	return meal, errorMapper(err)
}

func (uc MealUseCase) Create(ctx context.Context, date *time.Time) (entity.Meal, error) {
	var meal = &entity.Meal{
		IsAlone:     true,
		Freshness:   entity.Fresh,
		StressLevel: 0,
	}
	if date == nil || date.IsZero() {
		meal.Date = time.Now()
	} else {
		meal.Date = *date
	}

	err := uc.meals.CreateMeal(ctx, meal)
	return *meal, errorMapper(err)
}

func (uc MealUseCase) Delete(ctx context.Context, id uint) (ings entity.Ingredients, err error) {
	err = uc.meals.DeleteMeal(ctx, id)
	if err == nil {
		ings, err := uc.ingredients.ListIngredients(ctx)
		return ings, errorMapper(err)
	}
	return ings, errorMapper(err)
}

func (uc MealUseCase) Patch(ctx context.Context, id uint, date *time.Time, freshness *entity.Freshness, stressLevel *uint8, isAlone *bool) error {
	return errorMapper(uc.meals.PatchMeal(ctx, id, date, freshness, stressLevel, isAlone))
}
