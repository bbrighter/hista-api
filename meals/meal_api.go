package meals

import (
	"context"
	"sort"
	"time"

	"encore.dev/beta/errs"
)

type MealMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

type MealResponse struct {
	ID    uint           `json:"id"`
	Date  time.Time      `json:"date"`
	Foods []FoodResponse `json:"foods"`
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
		ID:    meal.ID,
		Date:  meal.Date,
		Foods: foodsToFoodsResponse(meal.Foods),
	}
	return resp
}

func mealsToMealsResponse(meals []Meal) MealsResponse {
	var resps []MealMetaResponse
	for _, m := range meals {
		resps = append(resps, mealToMealMetaResponse(m))
	}
	sort.Slice(resps, func(i, j int) bool {
		return resps[i].Date.Sub(resps[j].Date) > 0
	})
	return MealsResponse{resps}
}

// encore:api auth method=GET path=/meals
func (service Service) GetMeals(ctx context.Context) (MealsResponse, error) {
	meals := service.getMeals()
	return mealsToMealsResponse(meals), nil
}

type MealParams struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=POST path=/meals
func (service Service) PostMeal(ctx context.Context, params MealParams) (IDResponse, error) {
	id, err := service.createMeal(params.Date)
	if err != nil {
		return IDResponse{ID: id}, &errs.Error{Code: errs.InvalidArgument, Message: "IngredientID provided but not in database", Details: errs.Details(err)}
	}
	return IDResponse{ID: id}, err
}

// encore:api auth method=GET path=/meals/:id
func (service Service) GetMeal(ctx context.Context, id uint) (MealResponse, error) {
	meal, err := service.getMeal(id)
	return mealToMealResponse(meal), err
}

// encore:api auth method=DELETE path=/meals/:id
func (service Service) DeleteMeal(ctx context.Context, id uint) error {
	return service.deleteMeal(id)
}

// encore:api auth method=PATCH path=/meals/:id
func (service Service) PatchMealTime(ctx context.Context, id uint, params MealParams) error {
	return service.db.Model(&Meal{ID: id}).Update("date", params.Date).Error
}
