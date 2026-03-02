package internal

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
)

type (
	IFoodRepository interface {
		ListFoods(ctx context.Context, mealId uint) ([]*entity.Food, error)
		CreateFoodByName(ctx context.Context, food *entity.Food, ingredientName string) error
		CreateFoodByID(ctx context.Context, food *entity.Food) error
		DeleteFood(ctx context.Context, foodId uint) error
		GetFood(ctx context.Context, foodId uint) (entity.Food, error)
		UpdateFood(ctx context.Context, foodId uint, column string, value any) error
	}

	IFoodUseCase interface {
		List(ctx context.Context, mealId uint) (entity.Foods, error)
		Create(ctx context.Context, mealId uint, ingredientName string, ingredientId uint) (entity.Food, entity.Ingredients, error)
		Delete(ctx context.Context, foodId uint) (entity.Ingredients, error)
		ChangeCondition(ctx context.Context, foodId uint, newCond entity.FoodCondition) error
		ChangeAmount(ctx context.Context, foodId uint, amount *int) error
	}
)

type FoodUseCase struct {
	food        IFoodRepository
	ingredients IIngredientRepository
}

func NewFoodUseCase(foodRepo IFoodRepository, ingRepo IIngredientRepository) FoodUseCase {
	return FoodUseCase{food: foodRepo, ingredients: ingRepo}
}

func (uc FoodUseCase) List(ctx context.Context, mealId uint) (entity.Foods, error) {
	foods, err := uc.food.ListFoods(ctx, mealId)
	return foods, err
}

func (uc FoodUseCase) Create(ctx context.Context, mealId uint, ingredientName string, ingredientId uint) (entity.Food, entity.Ingredients, error) {
	var food = &entity.Food{MealID: mealId, Condition: entity.Cooked}
	var err error
	var ingredients entity.Ingredients
	if ingredientId != 0 {
		food.IngredientID = ingredientId
		err = uc.food.CreateFoodByID(ctx, food)
		if err != nil {
			return *food, ingredients, err
		}
		*food, err = uc.food.GetFood(ctx, food.ID)
	} else if ingredientName != "" {
		err = uc.food.CreateFoodByName(ctx, food, ingredientName)
	} else {
		err = errors.ErrorAttributeMustBeSet("ingredientName or ingredientId")
	}
	if err == nil {
		ingredients, err = uc.ingredients.ListIngredients(ctx)
	}
	return *food, ingredients, err
}

func (uc FoodUseCase) Delete(ctx context.Context, foodId uint) (ingredients entity.Ingredients, err error) {
	err = uc.food.DeleteFood(ctx, foodId)
	if err == nil {
		ingredients, err := uc.ingredients.ListIngredients(ctx)
		return ingredients, err
	}
	return ingredients, err
}

func (uc FoodUseCase) ChangeAmount(ctx context.Context, foodId uint, amount *int) error {
	return uc.food.UpdateFood(ctx, foodId, "amount", amount)
}

func (uc FoodUseCase) ChangeCondition(ctx context.Context, foodId uint, newCond entity.FoodCondition) error {
	return uc.food.UpdateFood(ctx, foodId, "condition", string(newCond))
}
