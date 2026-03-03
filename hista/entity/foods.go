package entity

import (
	"encore.dev/types/uuid"
)

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
	Amount         *int
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
	ID         uint      `gorm:"primaryKey"`
	PIID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name       string    `gorm:"uniqueIndex"`
	IsArchived bool      `gorm:"not null"`
}

func (i *Ingredient) SetPiid(id uuid.UUID) {
	i.PIID = id
}

type Ingredients []*Ingredient

type FoodResponse struct {
	ID         uint               `json:"id"`
	Ingredient IngredientResponse `json:"ingredient"`
	Condition  FoodCondition      `json:"foodCondition"`
	Amount     *int               `json:"amount,omitempty" encore:"optional"`
}

type FoodsResponse struct {
	Foods []FoodResponse `json:"foods"`
}

type IngredientResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"isArchived"`
}

type IngredientsResponse struct {
	Ingredients []IngredientResponse `json:"ingredients"`
}

func (i Ingredient) ToIngredientResponse() IngredientResponse {
	return IngredientResponse{
		ID:         i.ID,
		Name:       i.Name,
		IsArchived: i.IsArchived,
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
		Amount:     food.Amount,
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
