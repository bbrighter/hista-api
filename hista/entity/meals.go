package entity

import (
	"sort"
	"time"

	"encore.dev/types/uuid"
)

type Freshness uint8

const (
	Fresh   Freshness = 0
	SameDay Freshness = 1
	Older   Freshness = 2
)

type Meal struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	PIID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Date        time.Time
	Freshness   Freshness
	StressLevel uint8
	IsAlone     bool
	Foods       []Food `gorm:"constraint:OnDelete:CASCADE"`
}

func (m *Meal) SetPiid(id uuid.UUID) {
	m.PIID = id
}

type Meals []*Meal

type Food struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	PIID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Condition      FoodCondition
	Ingredient     Ingredient `gorm:"foreignKey:IngredientID,IngredientPIID;references:ID,PIID"`
	IngredientID   uint
	IngredientPIID uuid.UUID
	Meal           Meal `gorm:"foreignKey:MealID,MealPIID;references:ID,PIID"`
	MealID         uint
	MealPIID       uuid.UUID
}

func (f *Food) SetPiid(id uuid.UUID) {
	f.PIID = id
	f.IngredientPIID = id
	f.MealPIID = id
}

type Foods []*Food
type NonPtFoods []Food

type FoodCondition string

const (
	Raw    FoodCondition = "raw"
	Cooked FoodCondition = "cooked"
)

type Ingredient struct {
	ID   uint      `gorm:"primaryKey"`
	PIID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"uniqueIndex"`
}

func (i *Ingredient) SetPiid(id uuid.UUID) {
	i.PIID = id
}

type Ingredients []*Ingredient

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

func (i Ingredient) ToIngredientResponse() IngredientResponse {
	return IngredientResponse{
		ID:   i.ID,
		Name: i.Name,
	}
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

func (foods NonPtFoods) ToFoodsResponse() []FoodResponse {
	var foodsResponse = []FoodResponse{}
	for _, f := range foods {
		foodsResponse = append(foodsResponse, f.ToFoodResponse())
	}
	return foodsResponse
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
		Foods:       NonPtFoods(meal.Foods).ToFoodsResponse(),
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

type PostMealParams struct {
	Date time.Time `json:"date"`
}

type PatchMealParams struct {
	Date        *time.Time `json:"date" encore:"optional"`
	Freshness   *Freshness `json:"freshness" encore:"optional"`
	StressLevel *uint8     `json:"stressLevel" encore:"optional"`
	IsAlone     *bool      `json:"isAlone" encore:"optional"`
}
