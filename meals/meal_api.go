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
	ID    uint           `json:"id"`
	Date  time.Time      `json:"date"`
	Foods []FoodResponse `json:"foods"`
}

type MealsResponse struct {
	Meals []MealMetaResponse `json:"meals"`
}

// encore:api auth method=GET path=/meals
func (service Service) GetMeals(ctx context.Context) (MealsResponse, error) {
	meals, err := service.getMeals()
	return meals.toMealsResponse(), err
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
	return meal.toMealResponse(), err
}

// encore:api auth method=DELETE path=/meals/:id
func (service Service) DeleteMeal(ctx context.Context, id uint) error {
	return service.deleteMeal(id)
}

// encore:api auth method=PATCH path=/meals/:id
func (service Service) PatchMealTime(ctx context.Context, id uint, params MealParams) error {
	return service.db.Model(&Meal{ID: id}).Update("date", params.Date).Error
}
