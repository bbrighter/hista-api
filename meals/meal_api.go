package meals

import (
	"context"
	"time"
)

type MealMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

type MealResponse struct {
	ID          uint           `json:"id"`
	Date        time.Time      `json:"date"`
	Freshness   Freshness      `json:"freshness"`
	StressLevel uint8          `json:"stressLevel"`
	IsAlone     bool           `json:"isAlone"`
	Foods       []FoodResponse `json:"foods"`
}

type MealsResponse struct {
	Meals []MealMetaResponse `json:"meals"`
}

// encore:api auth method=GET path=/meals
func (service *Service) GetMeals(ctx context.Context) (MealsResponse, error) {
	var meals = new(Meals)
	var err error = meals.get(service)
	return meals.toMealsResponse(), err
}

type MealParams struct {
	Date        *time.Time `json:"date" encore:"optional"`
	Freshness   *Freshness `json:"freshness" encore:"optional"`
	StressLevel *uint8     `json:"stressLevel" encore:"optional"`
	IsAlone     *bool      `json:"isAlone" encore:"optional"`
}

// encore:api auth method=POST path=/meals
func (service *Service) PostMeal(ctx context.Context, params MealParams) (IDResponse, error) {
	var meal = Meal{Date: *params.Date, IsAlone: true}
	var err error = meal.create(service)
	return IDResponse{ID: meal.ID}, err
}

// encore:api auth method=GET path=/meals/:id
func (service *Service) GetMeal(ctx context.Context, id uint) (MealResponse, error) {
	var meal = Meal{ID: id}
	var err error = meal.get(service)
	return meal.toMealResponse(), err
}

// encore:api auth method=DELETE path=/meals/:id
func (service *Service) DeleteMeal(ctx context.Context, id uint) error {
	var meal = Meal{ID: id}
	return meal.delete(service)
}

// encore:api auth method=PATCH path=/meals/:id
func (service *Service) PatchMeal(ctx context.Context, id uint, params MealParams) error {
	var meal = Meal{ID: id}

	var patchParams = PatchParams(params)
	return meal.patch(service, patchParams)
}
