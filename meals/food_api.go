package meals

import (
	"context"

	"encore.app/errors"
)

type FoodResponse struct {
	ID         uint               `json:"id"`
	Ingredient IngredientResponse `json:"ingredient"`
	Condition  FoodCondition      `json:"foodCondition"`
}

type FoodsResponse struct {
	Foods []FoodResponse `json:"foods"`
}

// encore:api auth method=GET path=/meal/:mealId/foods
func (service *Service) GetFoods(ctx context.Context, mealId uint) (FoodsResponse, error) {
	var foods Foods = getFoods(service, mealId)
	var resp []FoodResponse = foods.toFoodsResponse()
	return FoodsResponse{Foods: resp}, nil
}

type FoodParams struct {
	IngredientName string        `json:"ingredientName" encore:"optional"`
	IngredientID   uint          `json:"ingredientId" encore:"optional"`
	Condition      FoodCondition `json:"condition" validate:"oneof=raw cooked"`
}

type PostFoodResponse struct {
	Food        FoodResponse        `json:"food"`
	Ingredients IngredientsResponse `json:"ingredients"`
}

// encore:api auth method=POST path=/meal/:mealId/foods
func (service *Service) PostFood(ctx context.Context, mealId uint, params FoodParams) (PostFoodResponse, error) {
	if params.IngredientID == 0 && params.IngredientName == "" {
		return PostFoodResponse{}, errors.ErrorAttributeMustBeSet("ingredientId or ingredientName")
	}
	var food = &Food{
		MealID:    mealId,
		Condition: params.Condition,
	}
	var ingredients Ingredients
	var err error
	if params.IngredientID != 0 {
		food.IngredientID = params.IngredientID
		ingredients, err = food.createByID(service)
	} else if params.IngredientName != "" {
		ingredients, err = food.createByName(service, params.IngredientName)
	}
	return PostFoodResponse{Food: food.toFoodResponse(), Ingredients: ingredients.toIngredientsResponse()}, err
}

// encore:api auth method=DELETE path=/meal/:mealId/foods/:foodId
func (service *Service) DeleteFood(ctx context.Context, mealId uint, foodId uint) (IngredientsResponse, error) {
	var food = Food{ID: foodId, MealID: mealId}
	var err error
	var ingredients Ingredients
	ingredients, err = food.delete(service)
	return ingredients.toIngredientsResponse(), err
}

type FoodConditionParams struct {
	Condition string `query:"condition"`
}

// encore:api auth method=PATCH path=/meal/:mealId/foods/:foodId/condition
func (service *Service) PatchFoodCondition(ctx context.Context, mealId uint, foodId uint, params FoodConditionParams) error {
	var food = &Food{ID: foodId}
	condition, err := stringToFoodCondition(params.Condition)
	if err != nil {
		return err
	}
	return food.changeCondition(service, condition)
}
