package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/internal/meals"
	"encore.dev/types/uuid"
)

type TemplateItemResponse struct {
	ID           uint   `json:"item"`
	IngredientId uint   `json:"ingredientId"`
	Condition    string `json:"condition"`
}

type TemplateResponse struct {
	ID    uint                   `json:"id"`
	Name  string                 `json:"name"`
	Items []TemplateItemResponse `json:"items"`
}

type TemplateListResponse struct {
	Templates []TemplateResponse `json:"templates"`
}

func toTemplateListResponse(mts meals.Templates) TemplateListResponse {
	var templates = []TemplateResponse{}
	for _, mt := range mts {
		var items = []TemplateItemResponse{}
		for _, i := range mt.Items {
			items = append(items, TemplateItemResponse{
				ID:           i.ID,
				IngredientId: i.IngredientID,
				Condition:    string(i.Condition),
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

// encore:api auth method=GET path=/piid/:piid/templates
func (service *Service) ListTemplates(ctx context.Context, piid uuid.UUID) (TemplateListResponse, error) {
	templates, err := service.meals.ListTemplatesAndItems(ctx)
	if err != nil {
		return TemplateListResponse{}, errors.MapError(err)
	}
	return toTemplateListResponse(templates), nil
}

type TemplateParams struct {
	Name  string               `json:"name"`
	Items []TemplateItemParams `json:"items"`
}

type TemplateItemParams struct {
	IngredientId uint   `json:"ingredientId"`
	Condition    string `json:"condition"`
}

func (p TemplateParams) ToItems() []meals.TemplateItem {
	var items = []meals.TemplateItem{}
	for _, i := range p.Items {
		items = append(items, meals.TemplateItem{
			IngredientID: i.IngredientId,
			Condition:    meals.FoodCondition(i.Condition),
		})
	}
	return items
}

type IDResponse struct {
	ID uint `json:"id"`
}

// encore:api auth method=POST path=/piid/:piid/templates
func (service *Service) PostTemplate(ctx context.Context, piid uuid.UUID, params TemplateParams) (IDResponse, error) {
	id, err := service.meals.CreateTemplate(ctx, params.Name, params.ToItems())
	return IDResponse{ID: id}, errors.MapError(err)
}

// encore:api auth method=PUT path=/piid/:piid/templates/:id
func (service *Service) PutTemplate(ctx context.Context, piid uuid.UUID, id uint, params TemplateParams) error {
	err := service.meals.ReplaceTemplate(ctx, id, params.Name, params.ToItems())
	return errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/templates/:id
func (service *Service) DeleteTemplate(ctx context.Context, piid uuid.UUID, id uint) error {
	err := service.meals.DeleteTemplate(ctx, id)
	return errors.MapError(err)
}
