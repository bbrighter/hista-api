package entity

import (
	"encore.dev/types/uuid"
)

type Ingredient struct {
	ID         uint      `gorm:"primaryKey"`
	PIID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name       string    `gorm:"uniqueIndex"`
	IsArchived bool      `gorm:"not null"`
	Nutrition  `gorm:"embeddedPrefix:nutrition_"`
}

func (i *Ingredient) SetPiid(id uuid.UUID) {
	i.PIID = id
}

type Ingredients []*Ingredient

type IngredientResponse struct {
	ID         uint           `json:"id"`
	Name       string         `json:"name"`
	IsArchived bool           `json:"isArchived"`
	Nutrition  *NutritionResp `json:"nutrition" encore:"optional"`
}

type IngredientsResponse struct {
	Ingredients []IngredientResponse `json:"ingredients"`
}

func (i Ingredient) ToIngredientResponse() IngredientResponse {
	return IngredientResponse{
		ID:         i.ID,
		Name:       i.Name,
		IsArchived: i.IsArchived,
		Nutrition:  i.toNutritionResp(),
	}
}

func (ingredients Ingredients) ToIngredientsResponse() IngredientsResponse {
	var resps []IngredientResponse
	for _, ing := range ingredients {
		resps = append(resps, ing.ToIngredientResponse())
	}
	return IngredientsResponse{resps}
}

type Nutrition struct {
	Protein      *int
	Carbohydrate *int
	Fat          *int
	Fiber        *int
}

type NutritionResp struct {
	Protein      int `json:"protein"`
	Carbohydrate int `json:"carbohydrate"`
	Fat          int `json:"fat"`
	Fiber        int `json:"fiber"`
}

func (n Nutrition) toNutritionResp() *NutritionResp {
	if n.Protein == nil || n.Fat == nil || n.Fiber == nil || n.Carbohydrate == nil {
		return nil
	}
	return &NutritionResp{
		Protein:      *n.Protein,
		Carbohydrate: *n.Carbohydrate,
		Fat:          *n.Fat,
		Fiber:        *n.Fiber,
	}
}

type PatchNutritionParams struct {
	Protein      int `json:"protein"`
	Carbohydrate int `json:"carbohydrate"`
	Fat          int `json:"fat"`
	Fiber        int `json:"fiber"`
}

func (p *PatchNutritionParams) ToNutrition() *Nutrition {
	if p == nil {
		return nil
	}
	return &Nutrition{
		Protein:      &p.Protein,
		Carbohydrate: &p.Carbohydrate,
		Fat:          &p.Fat,
		Fiber:        &p.Fiber,
	}
}
