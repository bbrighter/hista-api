package entity

import "encore.dev/types/uuid"

type Template struct {
	ID    uint           `gorm:"primaryKey"`
	PIID  uuid.UUID      `gorm:"type:uuid;primaryKey;uniqueIndex:idx_template_name_piid"`
	Name  string         `gorm:"uniqueIndex:idx_template_name_piid"`
	Items []TemplateItem `gorm:"constraint:OnDelete:CASCADE"`
}

func (mt *Template) SetPiid(id uuid.UUID) {
	mt.PIID = id

	for i := range mt.Items {
		mt.Items[i].SetPiid(id)
	}
}

type Templates []*Template

type TemplateItem struct {
	ID             uint      `gorm:"primaryKey"`
	PIID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Condition      FoodCondition
	TemplateID     uint
	TemplatePIID   uuid.UUID `gorm:"type:uuid"`
	IngredientID   uint
	IngredientPIID uuid.UUID `gorm:"type:uuid"`
}

func (mti *TemplateItem) SetPiid(id uuid.UUID) {
	mti.PIID = id
	mti.TemplatePIID = id
	mti.IngredientPIID = id
}

type TemplateItemResponse struct {
	ID           uint          `json:"item"`
	IngredientId uint          `json:"ingredientId"`
	Condition    FoodCondition `json:"condition"`
}

type TemplateResponse struct {
	ID    uint                   `json:"id"`
	Name  string                 `json:"name"`
	Items []TemplateItemResponse `json:"items"`
}

type TemplateListResponse struct {
	Templates []TemplateResponse `json:"templates"`
}

func (mts Templates) ToResponse() TemplateListResponse {
	var templates = []TemplateResponse{}
	for _, mt := range mts {
		var items = []TemplateItemResponse{}
		for _, i := range mt.Items {
			items = append(items, TemplateItemResponse{
				ID:           i.ID,
				IngredientId: i.IngredientID,
				Condition:    i.Condition,
			})
		}
		templates = append(templates, TemplateResponse{
			ID:    mt.ID,
			Name:  mt.Name,
			Items: items,
		})
	}
	return TemplateListResponse{Templates: templates}
}
