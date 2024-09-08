package entity

import (
	"sort"
	"time"
)

type Freshness uint8

const (
	Fresh   Freshness = 0
	SameDay Freshness = 1
	Older   Freshness = 2
)

type Meal struct {
	ID          uint
	Date        time.Time
	Freshness   Freshness
	StressLevel uint8
	IsAlone     bool
	Foods       []Food `gorm:"constraint:OnDelete:CASCADE"`
}

type Meals []Meal

type Food struct {
	ID           uint
	Ingredient   Ingredient
	IngredientID uint
	Condition    FoodCondition
	MealID       uint
}
type Foods []Food

type FoodCondition string

const (
	Raw    FoodCondition = "raw"
	Cooked FoodCondition = "cooked"
)

type Ingredient struct {
	ID   uint
	Name string `gorm:"uniqueIndex"`
}

type Ingredients []Ingredient

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

type FoodResponse struct {
	ID         uint               `json:"id"`
	Ingredient IngredientResponse `json:"ingredient"`
	Condition  FoodCondition      `json:"foodCondition"`
}

type FoodsResponse struct {
	Foods []FoodResponse `json:"foods"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

type IngredientResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type IngredientsResponse struct {
	Ingredients []IngredientResponse `json:"ingredients"`
}

func (ingredient Ingredient) ToIngredientResponse() IngredientResponse {
	return IngredientResponse(ingredient)
}

func (ingredients Ingredients) ToIngredientsResponse() IngredientsResponse {
	var resps []IngredientResponse
	for _, ing := range ingredients {
		resps = append(resps, ing.ToIngredientResponse())
	}
	return IngredientsResponse{resps}
}

func (food Food) ToFoodResponse() FoodResponse {
	return FoodResponse{
		ID:         food.ID,
		Ingredient: food.Ingredient.ToIngredientResponse(),
		Condition:  food.Condition,
	}
}

func (foods Foods) ToFoodsResponse() []FoodResponse {
	var foodsResponse = []FoodResponse{}
	for _, f := range foods {
		foodsResponse = append(foodsResponse, f.ToFoodResponse())
	}
	return foodsResponse
}

func (meal Meal) ToMealMetaResponse() MealMetaResponse {
	return MealMetaResponse{
		ID:   meal.ID,
		Date: meal.Date,
	}
}

func (meal Meal) ToMealResponse() MealResponse {
	var resp = MealResponse{
		ID:          meal.ID,
		Date:        meal.Date,
		Foods:       Foods(meal.Foods).ToFoodsResponse(),
		Freshness:   meal.Freshness,
		StressLevel: meal.StressLevel,
		IsAlone:     meal.IsAlone,
	}
	return resp
}

func (meals Meals) ToMealsResponse() MealsResponse {
	var resps []MealMetaResponse
	for _, m := range meals {
		resps = append(resps, m.ToMealMetaResponse())
	}
	sort.Slice(resps, func(i, j int) bool {
		return resps[i].Date.Sub(resps[j].Date) > 0
	})
	return MealsResponse{resps}
}

type MealParams struct {
	Date        *time.Time `json:"date" encore:"optional"`
	Freshness   *Freshness `json:"freshness" encore:"optional"`
	StressLevel *uint8     `json:"stressLevel" encore:"optional"`
	IsAlone     *bool      `json:"isAlone" encore:"optional"`
}

type PatchParams struct {
	Date        *time.Time
	Freshness   *Freshness
	StressLevel *uint8
	IsAlone     *bool
}
