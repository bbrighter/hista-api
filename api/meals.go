package api

import (
	"context"

	"encore.app/entity"
)

// encore:api auth method=GET path=/meals
func (service *Service) GetMeals(ctx context.Context) (entity.MealsResponse, error) {
	meals := service.mealUC.ListMeals()
	return meals.ToMealsResponse(), nil
}

// encore:api auth method=POST path=/meals
func (service *Service) PostMeal(ctx context.Context, params entity.MealParams) (entity.IDResponse, error) {
	id, err := service.mealUC.CreateMeal(params.Date)
	return entity.IDResponse{ID: id}, err
}

// encore:api auth method=GET path=/meals/:id
func (service *Service) GetMeal(ctx context.Context, id uint) (entity.MealResponse, error) {
	var meal = entity.Meal{ID: id}
	meal, err := service.mealUC.GetMeal(id)
	return meal.ToMealResponse(), err
}

// encore:api auth method=DELETE path=/meals/:id
func (service *Service) DeleteMeal(ctx context.Context, id uint) error {
	return service.mealUC.DeleteMeal(id)
}

// encore:api auth method=PATCH path=/meals/:id
func (service *Service) PatchMeal(ctx context.Context, id uint, params entity.MealParams) error {
	var patchParams = entity.PatchParams(params)
	return service.mealUC.PatchMeal(id, patchParams.Date, params.Freshness, params.StressLevel, params.IsAlone)
}

// encore:api auth method=GET path=/meal/:mealId/foods
func (service *Service) GetFoods(ctx context.Context, mealId uint) (entity.FoodsResponse, error) {
	var foods entity.Foods = service.food.ListFoods(mealId)
	var resp []entity.FoodResponse = foods.ToFoodsResponse()
	return entity.FoodsResponse{Foods: resp}, nil
}

type FoodParams struct {
	IngredientName string               `json:"ingredientName" encore:"optional"`
	IngredientID   uint                 `json:"ingredientId" encore:"optional"`
	Condition      entity.FoodCondition `json:"condition" validate:"oneof=raw cooked"`
}

type PostFoodResponse struct {
	Food        entity.FoodResponse        `json:"food"`
	Ingredients entity.IngredientsResponse `json:"ingredients"`
}

// encore:api auth method=POST path=/meal/:mealId/foods
func (service *Service) PostFood(ctx context.Context, mealId uint, params FoodParams) (PostFoodResponse, error) {
	food, ingredients, err := service.food.CreateFood(mealId, params.IngredientName, params.IngredientID)
	return PostFoodResponse{Food: food.ToFoodResponse(), Ingredients: ingredients.ToIngredientsResponse()}, err
}
