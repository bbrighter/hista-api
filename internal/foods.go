package internal

import (
	"encore.app/entity"
	"encore.app/errors"
)

type FoodUseCase struct {
	food        IFoodRepository
	ingredients IIngredientRepository
}

func NewFoodUseCase(foodRepo IFoodRepository, ingRepo IIngredientRepository) FoodUseCase {
	return FoodUseCase{food: foodRepo, ingredients: ingRepo}
}

func (uc FoodUseCase) ListFoods(mealId uint) entity.Foods {
	return uc.food.ListFoods(mealId)
}

func (uc FoodUseCase) CreateFood(mealId uint, ingredientName string, ingredientId uint) (food entity.Food, ingredients entity.Ingredients, err error) {
	var foodId uint
	if ingredientId != 0 {
		foodId, err = uc.food.CreateFoodByID(mealId, ingredientId)
	} else if ingredientName != "" {
		foodId, err = uc.food.CreateFoodByName(mealId, ingredientName)
	} else {
		err = errors.ErrorAttributeMustBeSet("ingredientName or ingredientId")
	}
	if err == nil {
		ingredients = uc.ingredients.ListIngredients()
		food = uc.food.GetFood(foodId)
	}
	return food, ingredients, err
}

func (uc FoodUseCase) DeleteFood(foodId uint) (ingredients entity.Ingredients, err error) {
	err = uc.food.DeleteFood(foodId)
	if err == nil {
		ingredients = uc.ingredients.ListIngredients()
	}
	return ingredients, err
}

func (uc FoodUseCase) ChangeCondition(foodId uint, newCond entity.FoodCondition) error {
	return uc.food.ChangeCondition(foodId, newCond)
}
