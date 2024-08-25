package api

import (
	"context"

	"encore.app/api/meals"
	"encore.app/entity"
	"encore.app/errors"
)

// encore:api auth method=GET path=/meals
func (service *Service) GetMeals(ctx context.Context) (entity.MealsResponse, error) {
	meals := service.mealRepository.List()
	return meals.ToMealsResponse(), nil
}

// encore:api auth method=POST path=/meals
func (service *Service) PostMeal(ctx context.Context, params entity.MealParams) (entity.IDResponse, error) {
	id, err := service.mealRepository.Create(*params.Date)
	return entity.IDResponse{ID: id}, err
}

// encore:api auth method=GET path=/meals/:id
func (service *Service) GetMeal(ctx context.Context, id uint) (entity.MealResponse, error) {
	var meal = entity.Meal{ID: id}
	meal, err := service.mealRepository.Get(id)
	return meal.ToMealResponse(), err
}

// encore:api auth method=DELETE path=/meals/:id
func (service *Service) DeleteMeal(ctx context.Context, id uint) error {
	return service.mealRepository.Delete(id)
}

// encore:api auth method=PATCH path=/meals/:id
func (service *Service) PatchMeal(ctx context.Context, id uint, params entity.MealParams) error {
	var patchParams = entity.PatchParams(params)
	return service.mealRepository.Patch(id, patchParams)
}

// encore:api auth method=GET path=/meal/:mealId/foods
func (service *Service) GetFoods(ctx context.Context, mealId uint) (entity.FoodsResponse, error) {
	var foods entity.Foods = service.mealRepository.ListFoods(mealId)
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
	if params.IngredientID == 0 && params.IngredientName == "" {
		return PostFoodResponse{}, errors.ErrorAttributeMustBeSet("ingredientId or ingredientName")
	}
	var food = &entity.Food{
		MealID:    mealId,
		Condition: params.Condition,
	}
	var err error
	if params.IngredientID != 0 {
		food.IngredientID = params.IngredientID
		err = service.mealRepository.CreateFoodByID(mealId, params.IngredientID)

	} else if params.IngredientName != "" {
		err = service.mealRepository.CreateFoodByName(mealId, params.IngredientName)
	}
	var ingredients entity.Ingredients = service.mealRepository.ListIngredients()
	return PostFoodResponse{Food: food.ToFoodResponse(), Ingredients: ingredients.ToIngredientsResponse()}, err
}

// encore:api auth method=DELETE path=/meal/:mealId/foods/:foodId
func (service *Service) DeleteFood(ctx context.Context, mealId uint, foodId uint) (entity.IngredientsResponse, error) {
	var err error = service.mealRepository.DeleteFood(foodId)
	var ingredients entity.Ingredients = service.mealRepository.ListIngredients()
	return ingredients.ToIngredientsResponse(), err
}

type FoodConditionParams struct {
	Condition string `query:"condition"`
}

// encore:api auth method=PATCH path=/meal/:mealId/foods/:foodId/condition
func (service *Service) PatchFoodCondition(ctx context.Context, mealId uint, foodId uint, params FoodConditionParams) error {
	var err error
	condition, err := meals.StringToFoodCondition(params.Condition)
	if err != nil {
		return err
	}
	return service.mealRepository.ChangeCondition(foodId, condition)

}
