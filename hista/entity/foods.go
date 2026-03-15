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

type FoodResponse struct {
	ID         uint               `json:"id"`
	Ingredient IngredientResponse `json:"ingredient"`
	Condition  FoodCondition      `json:"foodCondition"`
	Amount     *int               `json:"amount,omitempty" encore:"optional"`
}

type FoodsResponse struct {
	Foods []FoodResponse `json:"foods"`
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
