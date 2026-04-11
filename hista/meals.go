package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/meals
func (service *Service) ListMeals(ctx context.Context, piid uuid.UUID) (entity.MealsResponse, error) {
	meals, err := service.meals.List(ctx)
	return meals.ToMealsResponse(), errors.MapError(err)
}

// encore:api auth method=POST path=/piid/:piid/meals
func (service *Service) PostMeal(ctx context.Context, piid uuid.UUID, params entity.PostMealParams) (entity.MealResponse, error) {
	meal, err := service.meals.Create(ctx, params.Date)
	return meal.ToMealResponse(), errors.MapError(err)
}

// encore:api auth method=GET path=/piid/:piid/meals/:id
func (service *Service) GetMeal(ctx context.Context, piid uuid.UUID, id uint) (entity.MealResponse, error) {
	var meal = entity.Meal{ID: id}
	meal, err := service.meals.Get(ctx, id)
	return meal.ToMealResponse(), errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/meals/:id
func (service *Service) DeleteMeal(ctx context.Context, piid uuid.UUID, id uint) (entity.IngredientsResponse, error) {
	ings, err := service.meals.Delete(ctx, id)
	return ings.ToIngredientsResponse(), errors.MapError(err)
}

// encore:api auth method=PATCH path=/piid/:piid/meals/:id
func (service *Service) PatchMeal(ctx context.Context, piid uuid.UUID, id uint, params entity.PatchMealParams) error {
	return errors.MapError(service.meals.Patch(ctx, id, params.Date, params.Freshness, params.StressLevel, params.IsAlone))
}

// encore:api auth method=GET path=/piid/:piid/meal/:mealId/foods
func (service *Service) GetFoods(ctx context.Context, piid uuid.UUID, mealId uint) (entity.FoodsResponse, error) {
	foods, err := service.foods.List(ctx, mealId)
	return entity.FoodsResponse{Foods: foods.ToFoodsResponse()}, errors.MapError(err)
}

type FoodParams struct {
	IngredientName string `json:"ingredientName" encore:"optional"`
	IngredientID   uint   `json:"ingredientId" encore:"optional"`
}

type PostFoodResponse struct {
	Food        entity.FoodResponse        `json:"food"`
	Ingredients entity.IngredientsResponse `json:"ingredients"`
}

// encore:api auth method=POST path=/piid/:piid/meal/:mealId/foods
func (service *Service) PostFood(ctx context.Context, piid uuid.UUID, mealId uint, params FoodParams) (PostFoodResponse, error) {
	food, ingredients, err := service.foods.Create(ctx, mealId, params.IngredientName, params.IngredientID)
	return PostFoodResponse{Food: food.ToFoodResponse(), Ingredients: ingredients.ToIngredientsResponse()}, errors.MapError(err)
}

// encore:api method=POST path=/piid/:piid/meal/:mealId/foods/by-template/:templateId
func (service *Service) PostFoodByTemplate(ctx context.Context, piid uuid.UUID, mealId uint, templateId uint) (entity.FoodsResponse, error) {
	foods, err := service.mealTemplates.Apply(ctx, mealId, templateId)
	return entity.FoodsResponse{Foods: foods.ToFoodsResponse()}, errors.MapError(err)
}
