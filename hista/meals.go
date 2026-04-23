package hista

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/meals"
	"encore.dev/types/uuid"
)

type MealsResponse struct {
	Meals []MealMetaResponse `json:"meals"`
}

func toMealsResponse(meals meals.Meals) MealsResponse {
	var resps []MealMetaResponse
	for _, m := range meals {
		resps = append(resps, MealMetaResponse{
			ID:   m.ID,
			Date: m.Date,
		})
	}
	return MealsResponse{Meals: resps}
}

type MealMetaResponse struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

type MealResponse struct {
	MealMetaResponse
	Freshness     uint8 `json:"freshness"`
	StressLevel   uint8 `json:"stressLevel"`
	IsAlone       bool  `json:"isAlone"`
	FoodsResponse `json:"foods"`
}

func toMealResponse(meal meals.Meal) MealResponse {
	return MealResponse{
		MealMetaResponse: MealMetaResponse{
			ID:   meal.ID,
			Date: meal.Date,
		},
		Freshness:     uint8(meal.Freshness),
		StressLevel:   meal.StressLevel,
		IsAlone:       meal.IsAlone,
		FoodsResponse: toFoodsResponse(meal.Foods),
	}
}

// encore:api auth method=GET path=/piid/:piid/meals
func (service *Service) ListMeals(ctx context.Context, piid uuid.UUID) (MealsResponse, error) {
	meals, err := service.meals.ListMeals(ctx)
	if err != nil {
		return MealsResponse{}, errors.MapError(err)
	}
	return toMealsResponse(meals), nil
}

type PostMealParams struct {
	Date time.Time `json:"date"`
}

// encore:api auth method=POST path=/piid/:piid/meals
func (service *Service) PostMeal(ctx context.Context, piid uuid.UUID, params PostMealParams) (MealResponse, error) {
	meal, err := service.meals.CreateMeal(ctx, params.Date)
	if err != nil {
		return MealResponse{}, errors.MapError(err)
	}
	return toMealResponse(meal), nil
}

// encore:api auth method=GET path=/piid/:piid/meals/:id
func (service *Service) GetMeal(ctx context.Context, piid uuid.UUID, id uint) (MealResponse, error) {
	meal, err := service.meals.GetMeal(ctx, id)
	if err != nil {
		return MealResponse{}, errors.MapError(err)
	}
	return toMealResponse(meal), err
}

type IngredientsResponse struct {
	Ingredients []IngredientResponse `json:"ingredients"`
}

type IngredientResponse struct {
	ID         uint               `json:"id"`
	Name       string             `json:"name"`
	IsArchived bool               `json:"isArchived"`
	Nutrition  *NutritionResponse `json:"nutrition" encore:"optional"`
}

type NutritionResponse struct {
	Protein      float32 `json:"protein"`
	Carbohydrate float32 `json:"carbohydrate"`
	Fat          float32 `json:"fat"`
	Fiber        float32 `json:"fiber"`
}

func toNutritionResponse(n meals.Nutrition) *NutritionResponse {
	if n.Protein == nil || n.Fat == nil || n.Fiber == nil || n.Carbohydrate == nil {
		return nil
	}
	return &NutritionResponse{
		Protein:      *n.Protein,
		Carbohydrate: *n.Carbohydrate,
		Fat:          *n.Fat,
		Fiber:        *n.Fiber,
	}
}

func toIngredientsResponse(ingredients meals.Ingredients) IngredientsResponse {
	var resp = []IngredientResponse{}
	for _, i := range ingredients {
		resp = append(resp, IngredientResponse{
			ID:         i.ID,
			Name:       i.Name,
			IsArchived: i.IsArchived,
			Nutrition:  toNutritionResponse(i.Nutrition),
		})
	}
	return IngredientsResponse{Ingredients: resp}
}

// encore:api auth method=DELETE path=/piid/:piid/meals/:id
func (service *Service) DeleteMeal(ctx context.Context, piid uuid.UUID, id uint) (IngredientsResponse, error) {
	ings, err := service.meals.DeleteMeal(ctx, id)
	if err != nil {
		return IngredientsResponse{}, errors.MapError(err)
	}
	return toIngredientsResponse(ings), nil
}

type PatchMealParams struct {
	Date        *time.Time `json:"date" encore:"optional"`
	Freshness   *uint8     `json:"freshness" encore:"optional"`
	StressLevel *uint8     `json:"stressLevel" encore:"optional"`
	IsAlone     *bool      `json:"isAlone" encore:"optional"`
}

// encore:api auth method=PATCH path=/piid/:piid/meals/:id
func (service *Service) PatchMeal(ctx context.Context, piid uuid.UUID, id uint, params PatchMealParams) error {
	err := service.meals.UpdateMeal(ctx, id, params.Date, params.Freshness, params.StressLevel, params.IsAlone)
	return errors.MapError(err)
}

type FoodResponse struct {
	ID           uint   `json:"id"`
	IngredientId uint   `json:"ingredientId"`
	Condition    string `json:"foodCondition"`
	Amount       *int   `json:"amount,omitempty" encore:"optional"`
}

func toFoodResponse(f meals.Food) FoodResponse {
	return FoodResponse{
		ID:           f.ID,
		Amount:       f.Amount,
		Condition:    string(f.Condition),
		IngredientId: f.IngredientID,
	}
}

type FoodsResponse struct {
	Foods []FoodResponse `json:"foods"`
}

func toFoodsResponse(foods meals.Foods) FoodsResponse {
	var resps []FoodResponse
	for _, f := range foods {
		resps = append(resps, toFoodResponse(f))
	}
	return FoodsResponse{Foods: resps}
}

// encore:api auth method=GET path=/piid/:piid/meal/:mealId/foods
func (service *Service) GetFoods(ctx context.Context, piid uuid.UUID, mealId uint) (FoodsResponse, error) {
	foods, err := service.meals.ListFood(ctx, mealId)
	if err != nil {
		return FoodsResponse{}, errors.MapError(err)
	}
	return toFoodsResponse(foods), nil
}

type FoodParams struct {
	IngredientName string `json:"ingredientName" encore:"optional"`
	IngredientID   uint   `json:"ingredientId" encore:"optional"`
}

type PostFoodResponse struct {
	Food        FoodResponse        `json:"food"`
	Ingredients IngredientsResponse `json:"ingredients"`
}

// encore:api auth method=POST path=/piid/:piid/meal/:mealId/foods
func (service *Service) PostFood(ctx context.Context, piid uuid.UUID, mealId uint, params FoodParams) (PostFoodResponse, error) {
	var food meals.Food
	var err error
	if params.IngredientID != 0 {
		food, err = service.meals.CreateFoodById(ctx, mealId, params.IngredientID)
	} else {
		food, err = service.meals.CreateFoodByName(ctx, mealId, params.IngredientName)
	}
	if err != nil {
		return PostFoodResponse{}, errors.MapError(err)
	}
	ingredients, err := service.meals.ListIngredients(ctx)
	if err != nil {
		return PostFoodResponse{}, errors.MapError(err)
	}
	return PostFoodResponse{
		Food:        toFoodResponse(food),
		Ingredients: toIngredientsResponse(ingredients),
	}, nil
}

// encore:api auth method=POST path=/piid/:piid/meal/:mealId/foods/by-template/:templateId
func (service *Service) PostFoodByTemplate(ctx context.Context, piid uuid.UUID, mealId uint, templateId uint) (FoodsResponse, error) {
	foods, err := service.meals.ApplyTemplate(ctx, mealId, templateId)
	if err != nil {
		return FoodsResponse{}, errors.MapError(err)
	}
	return toFoodsResponse(foods), nil
}
