package entity

import (
	"encore.dev/types/uuid"
)

type Ingredient struct {
	ID         uint      `gorm:"primaryKey"`
	PIID       uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex:idx_name_piid"`
	Name       string    `gorm:"uniqueIndex:idx_name_piid"`
	IsArchived bool      `gorm:"not null"`
	Nutrition  `gorm:"embeddedPrefix:nutrition_"`
	Items      []TemplateItem `gorm:"constraint:OnDelete:CASCADE"`
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
	Protein      *float32
	Carbohydrate *float32
	Fat          *float32
	Fiber        *float32
}

type NutritionResp struct {
	Protein      float32 `json:"protein"`
	Carbohydrate float32 `json:"carbohydrate"`
	Fat          float32 `json:"fat"`
	Fiber        float32 `json:"fiber"`
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
	Protein      float32 `json:"protein"`
	Carbohydrate float32 `json:"carbohydrate"`
	Fat          float32 `json:"fat"`
	Fiber        float32 `json:"fiber"`
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
