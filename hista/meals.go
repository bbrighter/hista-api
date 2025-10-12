package hista

import (
	"context"

	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/meals
func (service *Service) ListMeals(ctx context.Context, piid uuid.UUID) (entity.MealsResponse, error) {
	meals, err := service.meals.List(ctx)
	return meals.ToMealsResponse(), err
}

// encore:api auth method=POST path=/piid/:piid/meals
func (service *Service) PostMeal(ctx context.Context, piid uuid.UUID, params entity.MealParams) (entity.MealResponse, error) {
	meal, err := service.meals.Create(ctx, params.Date)
	return meal.ToMealResponse(), err
}

// encore:api auth method=GET path=/piid/:piid/meals/:id
func (service *Service) GetMeal(ctx context.Context, piid uuid.UUID, id uint) (entity.MealResponse, error) {
	var meal = entity.Meal{ID: id}
	meal, err := service.meals.Get(ctx, id)
	return meal.ToMealResponse(), err
}

// encore:api auth method=DELETE path=/piid/:piid/meals/:id
func (service *Service) DeleteMeal(ctx context.Context, piid uuid.UUID, id uint) (entity.IngredientsResponse, error) {
	ings, err := service.meals.Delete(ctx, id)
	return ings.ToIngredientsResponse(), err
}

// encore:api auth method=PATCH path=/piid/:piid/meals/:id
func (service *Service) PatchMeal(ctx context.Context, piid uuid.UUID, id uint, params entity.MealParams) error {
	var patchParams = entity.PatchParams(params)
	return service.meals.Patch(ctx, id, patchParams.Date, params.Freshness, params.StressLevel, params.IsAlone)
}

// encore:api auth method=GET path=/piid/:piid/meal/:mealId/foods
func (service *Service) GetFoods(ctx context.Context, piid uuid.UUID, mealId uint) (entity.FoodsResponse, error) {
	foods, err := service.foods.List(ctx, mealId)
	return entity.FoodsResponse{Foods: foods.ToFoodsResponse()}, err
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

// encore:api auth method=POST path=/piid/:piid/meal/:mealId/foods
func (service *Service) PostFood(ctx context.Context, piid uuid.UUID, mealId uint, params FoodParams) (PostFoodResponse, error) {
	food, ingredients, err := service.foods.Create(ctx, mealId, params.IngredientName, params.IngredientID)
	return PostFoodResponse{Food: food.ToFoodResponse(), Ingredients: ingredients.ToIngredientsResponse()}, err
}
