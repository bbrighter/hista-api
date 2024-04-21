package meals

import (
	"context"
	"time"

	"encore.dev/beta/errs"
)

type MealMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

type MealResponse struct {
	MealMetaResponse
	Foods []FoodResponse `json:"foods"`
}

type FoodResponse struct {
	ID         uint               `json:"id"`
	Ingredient IngredientResponse `json:"ingredient"`
	Condition  FoodCondition      `json:"foodCondition"`
}

type MealsResponse struct {
	Meals []MealMetaResponse `json:"meals"`
}

func mealToMealMetaResponse(meal Meal) MealMetaResponse {
	return MealMetaResponse{
		ID:   meal.ID,
		Date: meal.Date,
	}
}

func mealToMealResponse(meal Meal) MealResponse {
	var foodsResponse []FoodResponse
	for _, f := range meal.Foods {
		foodsResponse = append(foodsResponse,
			FoodResponse{
				ID:         f.ID,
				Ingredient: ingredientToIngredientResponse(f.Ingredient),
				Condition:  f.Condition,
			})
	}
	var resp = MealResponse{
		MealMetaResponse: mealToMealMetaResponse(meal),
		Foods:            foodsResponse,
	}
	return resp
}

func mealsToMealsResponse(meals []Meal) MealsResponse {
	var resps []MealMetaResponse
	for _, m := range meals {
		resps = append(resps, mealToMealMetaResponse(m))
	}
	return MealsResponse{resps}
}

// encore:api public method=GET path=/meals
func (service Service) GetMeals(ctx context.Context) (MealsResponse, error) {
	meals := service.getMeals()
	return mealsToMealsResponse(meals), nil
}

type FoodParams struct {
	IngredientID uint          `json:"ingredientId"`
	Condition    FoodCondition `json:"condition"`
}

type MealParams struct {
	Date  time.Time    `json:"date"`
	Foods []FoodParams `json:"foods"`
}

// encore:api public method=POST path=/meals
func (service Service) PostMeal(ctx context.Context, params MealParams) (IDResponse, error) {
	var foods []Food
	for _, f := range params.Foods {
		foods = append(foods, Food{
			IngredientID: f.IngredientID,
			Condition:    f.Condition,
		})
	}
	id, err := service.createMeal(foods, params.Date)
	if err != nil {
		return IDResponse{ID: id}, &errs.Error{Code: errs.InvalidArgument, Message: "IngredientID provided but not in database", Details: errs.Details(err)}
	}
	return IDResponse{ID: id}, err
}

// encore:api public method=GET path=/meals/:id
func (service Service) GetMeal(ctx context.Context, id uint) (MealResponse, error) {
	meal, err := service.getMeal(id)
	return mealToMealResponse(meal), err
}
