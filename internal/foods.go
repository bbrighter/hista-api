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

func (uc FoodUseCase) List(mealId uint) entity.Foods {
	return uc.food.ListFoods(mealId)
}

func (uc FoodUseCase) Create(mealId uint, ingredientName string, ingredientId uint) (entity.Food, entity.Ingredients, error) {
	var food = &entity.Food{MealID: mealId, Condition: entity.Cooked}
	var err error
	var ingredients entity.Ingredients
	if ingredientId != 0 {
		food.IngredientID = ingredientId
		err = uc.food.CreateFoodByID(food)
		*food = uc.food.GetFood(food.ID)
	} else if ingredientName != "" {
		err = uc.food.CreateFoodByName(food, ingredientName)
	} else {
		err = errors.ErrorAttributeMustBeSet("ingredientName or ingredientId")
	}
	if err == nil {
		ingredients = uc.ingredients.ListIngredients()
	}
	return *food, ingredients, err
}

func (uc FoodUseCase) Delete(foodId uint) (ingredients entity.Ingredients, err error) {
	err = uc.food.DeleteFood(foodId)
	if err == nil {
		ingredients = uc.ingredients.ListIngredients()
	}
	return ingredients, err
}

func (uc FoodUseCase) ChangeCondition(foodId uint, newCond entity.FoodCondition) error {
	return uc.food.ChangeCondition(foodId, newCond)
}
